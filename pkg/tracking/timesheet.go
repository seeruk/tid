package tracking

import (
	"github.com/SeerUK/tid/proto"
)

// Timesheet wraps a ProtoBuf-generated proto.TimeSheet message with helper methods. No state should
// be kept in this struct.
type Timesheet struct {
	Message *proto.TrackingTimesheet
}

// NewTimesheet create a new instance of TimeSheet.
func NewTimesheet(message *proto.TrackingTimesheet) *Timesheet {
	return &Timesheet{
		Message: message,
	}
}

// AppendEntry appends a reference to an entry to the timesheet.
func (t *Timesheet) AppendEntry(entry *Entry) {
	t.Message.Entries = append(t.Message.Entries, entry.ShortKey())
}

// Key returns the key of the underlying message.
func (t *Timesheet) Key() string {
	return t.Message.Key
}
