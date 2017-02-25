package tracking

import "github.com/SeerUK/tid/proto"

// Status decorates a ProtoBuf-generated proto.Status type with helper methods.
type Status struct {
	*proto.Status
}

// NewStatus creates a new instance of StatusWrapper.
func NewStatus(proto *proto.Status) *Status {
	return &Status{
		proto,
	}
}

// IsActive checks if any timesheet is active.
func (s *Status) IsActive() bool {
	return s.State == proto.Status_STARTED || s.State == proto.Status_PAUSED
}
