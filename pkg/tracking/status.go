package tracking

import (
	"github.com/SeerUK/tid/proto"
)

// Status wraps a ProtoBuf-generated proto.TrackingStatus message with helper methods. No state
// should be kept in this struct.
type Status struct {
	Message *proto.TrackingStatus
}

// NewStatus create a new instance of Status.
func NewStatus() *Status {
	return &Status{
		Message: &proto.TrackingStatus{},
	}
}

// IsActive checks if any timesheet is active.
func (s *Status) IsActive() bool {
	return s.Message.GetState() == proto.TrackingStatus_STARTED
}

// Start updates the status to reflect the fact that a new entry is being tracked.
func (s *Status) Start(sheet *Timesheet, entry *Entry) {
	s.Message.State = proto.TrackingStatus_STARTED
	s.Message.Ref = &proto.TrackingStatusRef{
		Timesheet: sheet.Key(),
		Entry:     entry.Key(),
	}
}

// Stop updates the status to reflect that tracking has ended (at least temporarily).
func (s *Status) Stop() {
	s.Message.State = proto.TrackingStatus_STOPPED
}

// Ref gets the proto.TrackingStatusRef of the underlying ProtoBuf message.
func (s *Status) Ref() *proto.TrackingStatusRef {
	return s.Message.Ref
}
