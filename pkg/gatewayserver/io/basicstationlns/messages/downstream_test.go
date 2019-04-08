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

package messages

import (
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/basicstation"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func eui64Ptr(eui types.EUI64) *types.EUI64 { return &eui }
func timePtr(time time.Time) *time.Time     { return &time }
func TestDownlinkMessage(t *testing.T) {
	for _, tc := range []struct {
		Name                    string
		NSDownlinkMessage       ttnpb.DownlinkMessage
		GatewayIDs              ttnpb.GatewayIdentifiers
		ExpectedDownlinkMessage DownlinkMessage
	}{
		{
			Name: "SampleDownlink",
			NSDownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
					DevEUI:   eui64Ptr(types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}),
				},
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     868500000,
						Downlink: &ttnpb.TxSettings_Downlink{
							AntennaIndex: 2,
						},
						Timestamp: 1554300787,
					},
				},
			},
			GatewayIDs: ttnpb.GatewayIdentifiers{GatewayID: "test-gateway"},
			ExpectedDownlinkMessage: DownlinkMessage{
				DevEUI:      basicstation.EUI{EUI64: types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
				DeviceClass: 0,
				Pdu:         "Ymxhamthc25kJ3M==",
				RxDelay:     1,
				Rx2DR:       2,
				Rx2Freq:     868500000,
				RCtx:        2,
				Priority:    25,
				XTime:       1554300667,
			},
		},
		{
			Name: "WithAbsoluteTime",
			NSDownlinkMessage: ttnpb.DownlinkMessage{
				RawPayload: []byte("Ymxhamthc25kJ3M=="),
				EndDeviceIDs: &ttnpb.EndDeviceIdentifiers{
					DeviceID: "testdevice",
					DevEUI:   eui64Ptr(types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}),
				},
				Settings: &ttnpb.DownlinkMessage_Scheduled{
					Scheduled: &ttnpb.TxSettings{
						DataRateIndex: 2,
						Frequency:     869525000,
						Downlink: &ttnpb.TxSettings_Downlink{
							AntennaIndex: 2,
						},
						Time: timePtr(time.Unix(1554300787, 0)),
					},
				},
			},
			GatewayIDs: ttnpb.GatewayIdentifiers{GatewayID: "test-gateway"},
			ExpectedDownlinkMessage: DownlinkMessage{
				DevEUI:      basicstation.EUI{EUI64: types.EUI64{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
				DeviceClass: 1,
				Pdu:         "Ymxhamthc25kJ3M==",
				RxDelay:     1,
				Rx2DR:       2,
				Rx2Freq:     869525000,
				RCtx:        2,
				Priority:    25,
				GpsTime:     1554300787,
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			a := assertions.New(t)
			dnmsg := DownlinkMessage{}
			dnmsg.FromNSDownlinkMessage(tc.GatewayIDs, tc.NSDownlinkMessage, 0)
			if !a.So(dnmsg, should.Resemble, tc.ExpectedDownlinkMessage) {
				t.Fatalf("Invalid DownlinkMessage: %v", dnmsg)
			}
		})
	}
}
