package timesheet

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// IsActive checks if any timesheet is active by using a given status.
func IsActive(status proto.Status) bool {
	return status.State == proto.Status_STARTED || status.State == proto.Status_PAUSED
}

// UpdateEntryDuration sets the duration of a given timesheet's entry at a given index.
func UpdateEntryDuration(sheet *proto.TimeSheet, index int64) {
	// @todo: This method will make more sense when pausing is implemented. This is because the
	// duration will be set from the most recent start time, so it may be called more than once per
	// timesheet entry.
	entry := sheet.Entries[index]
	entry.Duration = uint64(time.Now().Unix()) - entry.StartTime
}
