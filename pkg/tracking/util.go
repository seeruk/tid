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

// createEntry creates a new timesheet entry instance.
func createEntry(duration uint64, startTime uint64, note string) *proto.TimeSheetEntry {
	return &proto.TimeSheetEntry{
		Duration:  duration,
		StartTime: startTime,
		Note:      note,
	}
}
