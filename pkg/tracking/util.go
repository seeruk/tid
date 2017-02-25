package tracking

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// AppendNewEntry appends a new timesheet entry to a given timesheet.
func AppendNewEntry(sheet *proto.TimeSheet, note string) {
	now := time.Now().Unix()

	sheet.Entries = append(sheet.Entries, createEntry(
		uint64(0),
		uint64(now),
		note,
	))
}

// UpdateEntryDuration sets the duration of a given timesheet's entry at a given index.
func UpdateEntryDuration(sheet *proto.TimeSheet, index int64) {
	// @todo: This method will make more sense when pausing is implemented. This is because the
	// duration will be set from the most recent start time, so it may be called more than once per
	// timesheet entry.
	entry := sheet.Entries[index]
	entry.Duration = uint64(time.Now().Unix()) - entry.StartTime
}

// @todo: Not happy with this... Maybe a "status" struct of some kind with a method on. That would
// make calling code much nicer looking. `status.StartEntry(&sheet)` or something.

// UpdateStatusStartEntry updates the given status to start with the latest entry on the given
// timesheet.
func UpdateStatusStartEntry(status *Status, sheet *proto.TimeSheet) {
	now := time.Now().Local()

	status.State = proto.Status_STARTED
	status.TimeSheetEntry = &proto.TimeSheetEntryRef{
		Date:  now.Format(KeyTimesheetFmt),
		Index: int64(len(sheet.Entries) - 1),
	}
}

// createEntry creates a new timesheet entry instance.
func createEntry(duration uint64, startTime uint64, note string) *proto.TimeSheetEntry {
	return &proto.TimeSheetEntry{
		Duration:  duration,
		StartTime: startTime,
		Note:      note,
	}
}
