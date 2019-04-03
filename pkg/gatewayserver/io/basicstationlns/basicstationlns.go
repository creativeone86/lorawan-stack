// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package basicstationlns

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	echo "github.com/labstack/echo/v4"
	"go.thethings.network/lorawan-stack/pkg/auth/rights"
	"go.thethings.network/lorawan-stack/pkg/basicstation"
	"go.thethings.network/lorawan-stack/pkg/errors"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io/basicstationlns/messages"
	"go.thethings.network/lorawan-stack/pkg/log"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/unique"
	"go.thethings.network/lorawan-stack/pkg/web"
)

const tokenExpiration = 3 * time.Minute

var (
	errEmptyGatewayEUI           = errors.Define("empty_gateway_eui", "empty gateway EUI")
	errMessageTypeNotImplemented = errors.DefineUnimplemented("message_type_not_implemented", "message of type `{type}` is not implemented")
)

// downlinkInfo is the information associated with a particular downlink
type downlinkInfo struct {
	correlationIDs []string
	txTime         time.Time
}

type srv struct {
	ctx      context.Context
	server   io.Server
	upgrader *websocket.Upgrader
	// token is the unique token associated with each Downlink.
	// It's passed to the BasicStation as the `diid` field and is returned as-is in the TxConfirmation if the downlink packet was put on air.
	// This is a free-running counter that is allowed to overflow and is cleaned up periodically by the garbage collector.
	token        int64
	correlations sync.Map
}

func (*srv) Protocol() string { return "basicstation" }

// New returns a new Basic Station frontend that can be registered in the web server.
func New(ctx context.Context, server io.Server) web.Registerer {
	ctx = log.NewContextWithField(ctx, "namespace", "gatewayserver/io/basicstation")
	s := &srv{
		ctx:      ctx,
		server:   server,
		upgrader: &websocket.Upgrader{},
	}
	go s.gc()
	return s
}

func (s *srv) RegisterRoutes(server *web.Server) {
	group := server.Group(ttnpb.HTTPAPIPrefix + "/gs/io/basicstation")
	group.GET("/discover", s.handleDiscover)
	group.GET("/traffic/:uid", s.handleTraffic)
}

func (s *srv) handleDiscover(c echo.Context) error {
	logger := log.FromContext(s.ctx).WithFields(log.Fields(
		"endpoint", "discover",
		"remote_addr", c.Request().RemoteAddr,
	))
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Debug("Failed to upgrade request to websocket connection")
		return err
	}
	defer ws.Close()

	_, data, err := ws.ReadMessage()
	if err != nil {
		logger.WithError(err).Debug("Failed to read message")
		return err
	}
	var req messages.DiscoverQuery
	if err := json.Unmarshal(data, &req); err != nil {
		logger.WithError(err).Debug("Failed to parse discover query message")
		return err
	}

	if req.EUI.IsZero() {
		writeDiscoverError(s.ctx, ws, "Invalid request")
		return errEmptyGatewayEUI
	}

	ids := ttnpb.GatewayIdentifiers{
		EUI: &req.EUI.EUI64,
	}
	ctx, ids, err := s.server.FillGatewayContext(s.ctx, ids)
	if err != nil {
		logger.WithError(err).Debug("Failed to fill gateway context")
		writeDiscoverError(ctx, ws, "Router not provisioned")
		return err
	}
	uid := unique.ID(ctx, ids)
	ctx = log.NewContextWithField(ctx, "gateway_uid", uid)

	scheme := "ws"
	if c.IsTLS() {
		scheme = "wss"
	}
	res := messages.DiscoverResponse{
		EUI: req.EUI,
		Muxs: basicstation.EUI{
			Prefix: "muxs",
		},
		URI: fmt.Sprintf("%s://%s%s", scheme, c.Request().Host, c.Echo().URI(s.handleTraffic, uid)),
	}
	data, err = json.Marshal(res)
	if err != nil {
		logger.WithError(err).Warn("Failed to marshal response message")
		writeDiscoverError(ctx, ws, "Router not provisioned")
		return err
	}
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		logger.WithError(err).Warn("Failed to write discover response message")
		return err
	}
	logger.Debug("Sent discover response message")
	return nil
}

func (s *srv) handleTraffic(c echo.Context) error {
	uid := c.Param("uid")
	ids, err := unique.ToGatewayID(uid)
	if err != nil {
		return err
	}
	ctx, err := unique.WithContext(s.ctx, uid)
	if err != nil {
		return err
	}
	ctx = log.NewContextWithField(s.ctx, "gateway_uid", uid)
	logger := log.FromContext(ctx).WithFields(log.Fields(
		"endpoint", "traffic",
		"remote_addr", c.Request().RemoteAddr,
	))
	fp, err := s.server.GetFrequencyPlan(ctx, ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to get frequency plan")
		return err
	}
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.WithError(err).Debug("Failed to upgrade request to websocket connection")
		return err
	}
	defer ws.Close()

	ctx = rights.NewContext(ctx, rights.Rights{
		GatewayRights: map[string]*ttnpb.Rights{
			uid: {
				Rights: []ttnpb.Right{ttnpb.RIGHT_GATEWAY_LINK},
			},
		},
	})
	conn, err := s.server.Connect(ctx, s, ids)
	if err != nil {
		logger.WithError(err).Warn("Failed to connect")
		return err
	}
	if err := s.server.ClaimDownlink(ctx, ids); err != nil {
		logger.WithError(err).Error("Failed to claim downlink")
		return err
	}
	defer func() {
		if err := s.server.UnclaimDownlink(ctx, ids); err != nil {
			logger.WithError(err).Error("Failed to unclaim downlink")
		}
	}()

	// Process downlinks in a separate go routine
	go func() {
		for {
			select {
			case <-conn.Context().Done():
				return
			case down := <-conn.Down():
				s.createNextToken()

				dnmsg := messages.DownlinkMessage{}
				dnmsg.FromNSDownlinkMessage(ids, *down, s.token)
				s.correlations.Store(s.token, downlinkInfo{
					correlationIDs: down.GetCorrelationIDs(),
					txTime:         time.Now(),
				})

				msg, err := dnmsg.MarshalJSON()
				if err != nil {
					logger.WithError(err).Error("Failed to marshal downlink message")
					continue
				}

				logger.Info("Sending downlink message")
				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					logger.WithError(err).Error("Failed to send downlink message")
				}
			}
		}
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			logger.WithError(err).Debug("Failed to read message")
			return err
		}
		typ, err := messages.Type(data)
		if err != nil {
			logger.WithError(err).Debug("Failed to parse message type")
			return err
		}
		logger = logger.WithFields(log.Fields(
			"upstream_type", typ,
		))

		switch typ {
		case messages.TypeUpstreamVersion:
			var version messages.Version
			if err := json.Unmarshal(data, &version); err != nil {
				logger.WithError(err).Debug("Failed to unmarshal version message")
				return err
			}
			logger = logger.WithFields(log.Fields(
				"station", version.Station,
				"firmware", version.Firmware,
				"model", version.Model,
			))
			cfg, err := messages.GetRouterConfig(*fp, version.IsProduction())
			if err != nil {
				logger.WithError(err).Warn("Failed to generate router configuration")
				return err
			}
			data, err = json.Marshal(cfg)
			if err != nil {
				logger.WithError(err).Warn("Failed to marshal response message")
				return err
			}
			if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.WithError(err).Warn("Failed to send router configuration")
				return err
			}

		case messages.TypeUpstreamJoinRequest:
			var jreq messages.JoinRequest
			if err := json.Unmarshal(data, &jreq); err != nil {
				logger.WithError(err).Debug("Failed to unmarshal join-request message")
				return err
			}
			up, err := jreq.ToUplinkMessage(ids, fp.BandID)
			if err != nil {
				logger.WithError(err).Debug("Failed to parse join-request message")
				return err
			}
			if err := conn.HandleUp(&up); err != nil {
				logger.WithError(err).Warn("Failed to handle uplink message")
			}

		case messages.TypeUpstreamUplinkDataFrame:
			var updf messages.UplinkDataFrame
			if err := json.Unmarshal(data, &updf); err != nil {
				logger.WithError(err).Debug("Failed to unmarshal uplink data frame")
				return err
			}
			up, err := updf.ToUplinkMessage(ids, fp.BandID)
			if err != nil {
				logger.WithError(err).Debug("Failed to parse uplink data frame")
				return err
			}
			if err := conn.HandleUp(&up); err != nil {
				logger.WithError(err).Warn("Failed to handle uplink message")
			}

		case messages.TypeUpstreamTxConfirmation:
			var txConf messages.TxConfirmation
			if err := json.Unmarshal(data, &txConf); err != nil {
				logger.WithError(err).Debug("Failed to unmarshal tx acknowledgement frame")
				return err
			}
			if value, ok := s.correlations.Load(txConf.Diid); ok {
				txAck := messages.ToTxAcknowledgment(value.(downlinkInfo).correlationIDs)
				if err := conn.HandleTxAck(&txAck); err != nil {
					logger.WithError(err).Warn("Failed to handle uplink message")
				}
				s.correlations.Delete(txConf.Diid)
			} else {
				logger.Debug("TxAck does not correspond to a downlink message or is received too late")
			}
		case messages.TypeUpstreamProprietaryDataFrame:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)
		case messages.TypeUpstreamRemoteShell:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)
		case messages.TypeUpstreamTimeSync:
			return errMessageTypeNotImplemented.WithAttributes("type", typ)

		default:
			// Unknown message types are ignored by the server
			logger.WithField("message_type", typ).Debug("Unknown message type")
		}
	}
}

// writeDiscoverError sends the error messages during the discovery on the WS connection to the station.
func writeDiscoverError(ctx context.Context, ws *websocket.Conn, msg string) {
	logger := log.FromContext(ctx)
	errMsg, err := json.Marshal(messages.DiscoverResponse{Error: msg})
	if err != nil {
		logger.WithError(err).Debug("Failed to marshal error message")
		return
	}
	if err := ws.WriteMessage(websocket.TextMessage, errMsg); err != nil {
		logger.WithError(err).Debug("Failed to write error response message")
	}
}

// createNextToken atomically increments the token value.
func (s *srv) createNextToken() {
	atomic.AddInt64(&s.token, 1)
}

// gc is the garbage collector that removes old items from the correlations map.
func (s *srv) gc() {
	gcTicker := time.NewTicker(tokenExpiration)
	for {
		select {
		case <-s.ctx.Done():
			gcTicker.Stop()
			return
		case <-gcTicker.C:
			s.correlations.Range(func(key interface{}, value interface{}) bool {
				if value.(downlinkInfo).txTime.Before(time.Now().Add(-tokenExpiration)) {
					s.correlations.Delete(key)
				}
				return true
			})
		}
	}
}