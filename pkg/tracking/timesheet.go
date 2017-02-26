package tracking

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// TimeSheet wraps a ProtoBuf-generated proto.TimeSheet message with helper methods. No state should
// be kept in this struct.
type TimeSheet struct {
	Message *proto.TimeSheet
}

// NewTimeSheet create a new instance of TimeSheet.
func NewTimeSheet(message *proto.TimeSheet) *TimeSheet {
	return &TimeSheet{
		Message: message,
	}
}

// AppendNewEntry appends a new timesheet entry to a given timesheet.
func (t *TimeSheet) AppendNewEntry(note string) proto.TimeSheetEntryRef {
	now := time.Now()

	t.Message.Entries = append(t.Message.Entries, createEntry(
		uint64(0),
		uint64(now.Unix()),
		note,
	))

	return proto.TimeSheetEntryRef{
		Date:  now.Format(KeyTimeSheetFmt),
		Index: int64(len(t.Message.Entries) - 1),
	}
}

// UpdateEntryDuration updates the duration of the entry referenced in the status.
func (t *TimeSheet) UpdateEntryDuration(status *Status) {
	// @todo: This method will make more sense when pausing is implemented. This is because the
	// duration will be set from the most recent start time, so it may be called more than once per
	// timesheet entry.
	entry := t.Message.Entries[status.TimeSheetEntry().Index]
	entry.Duration = uint64(time.Now().Unix()) - entry.StartTime
}

// createEntry creates a new proto.TimeSheetEntry instance.
func createEntry(duration uint64, startTime uint64, note string) *proto.TimeSheetEntry {
	return &proto.TimeSheetEntry{
		Duration:  duration,
		StartTime: startTime,
		Note:      note,
	}
}
