package tracking

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// Status wraps a ProtoBuf-generated proto.Status message with helper methods.
type Status struct {
	Message *proto.Status
}

// NewStatus creates a new instance of StatusWrapper.
func NewStatus(message *proto.Status) *Status {
	return &Status{
		Message: message,
	}
}

// IsActive checks if any timesheet is active.
func (s *Status) IsActive() bool {
	state := s.Message.GetState()

	return state == proto.Status_STARTED || state == proto.Status_PAUSED
}

// Start updates the status to reflect the fact that a new entry is being tracked.
func (s *Status) Start(timesheet *proto.TimeSheet) {
	now := time.Now().Local()

	s.Message.State = proto.Status_STARTED
	s.Message.TimeSheetEntry = &proto.TimeSheetEntryRef{
		Date:  now.Format(KeyTimeSheetFmt),
		Index: int64(len(timesheet.Entries) - 1),
	}
}

// TimeSheetEntry gets the proto.TimeSheetEntryRef in the underlying ProtoBuf message.
func (s *Status) TimeSheetEntry() *proto.TimeSheetEntryRef {
	return s.Message.TimeSheetEntry
}
