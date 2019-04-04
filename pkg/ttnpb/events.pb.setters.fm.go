// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	fmt "fmt"
	time "time"
)

func (dst *Event) SetFields(src *Event, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "name":
			if len(subs) > 0 {
				return fmt.Errorf("'name' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Name = src.Name
			} else {
				var zero string
				dst.Name = zero
			}
		case "time":
			if len(subs) > 0 {
				return fmt.Errorf("'time' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Time = src.Time
			} else {
				var zero time.Time
				dst.Time = zero
			}
		case "identifiers":
			if len(subs) > 0 {
				return fmt.Errorf("'identifiers' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Identifiers = src.Identifiers
			} else {
				dst.Identifiers = nil
			}
		case "data":
			if len(subs) > 0 {
				return fmt.Errorf("'data' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Data = src.Data
			} else {
				dst.Data = nil
			}
		case "correlation_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'correlation_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.CorrelationIDs = src.CorrelationIDs
			} else {
				dst.CorrelationIDs = nil
			}
		case "origin":
			if len(subs) > 0 {
				return fmt.Errorf("'origin' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Origin = src.Origin
			} else {
				var zero string
				dst.Origin = zero
			}
		case "context":
			if len(subs) > 0 {
				return fmt.Errorf("'context' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Context = src.Context
			} else {
				dst.Context = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamEventsRequest) SetFields(src *StreamEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "identifiers":
			if len(subs) > 0 {
				return fmt.Errorf("'identifiers' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Identifiers = src.Identifiers
			} else {
				dst.Identifiers = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamApplicationEventsRequest) SetFields(src *StreamApplicationEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'application_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ApplicationIDs = src.ApplicationIDs
			} else {
				dst.ApplicationIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamClientEventsRequest) SetFields(src *StreamClientEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "client_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'client_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ClientIDs = src.ClientIDs
			} else {
				dst.ClientIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamEndDeviceEventsRequest) SetFields(src *StreamEndDeviceEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "application_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'application_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.ApplicationIDs = src.ApplicationIDs
			} else {
				var zero string
				dst.ApplicationIDs = zero
			}
		case "device_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'device_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.DeviceIDs = src.DeviceIDs
			} else {
				dst.DeviceIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamGatewayEventsRequest) SetFields(src *StreamGatewayEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "gateway_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'gateway_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.GatewayIDs = src.GatewayIDs
			} else {
				dst.GatewayIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamOrganizationEventsRequest) SetFields(src *StreamOrganizationEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "organization_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'organization_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.OrganizationIDs = src.OrganizationIDs
			} else {
				dst.OrganizationIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}

func (dst *StreamUserEventsRequest) SetFields(src *StreamUserEventsRequest, paths ...string) error {
	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		switch name {
		case "user_ids":
			if len(subs) > 0 {
				return fmt.Errorf("'user_ids' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.UserIDs = src.UserIDs
			} else {
				dst.UserIDs = nil
			}
		case "tail":
			if len(subs) > 0 {
				return fmt.Errorf("'tail' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.Tail = src.Tail
			} else {
				var zero uint32
				dst.Tail = zero
			}
		case "after":
			if len(subs) > 0 {
				return fmt.Errorf("'after' has no subfields, but %s were specified", subs)
			}
			if src != nil {
				dst.After = src.After
			} else {
				dst.After = nil
			}

		default:
			return fmt.Errorf("invalid field: '%s'", name)
		}
	}
	return nil
}
