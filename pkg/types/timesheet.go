package types

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// TimesheetKeyDateFmt is the date formatting string for timesheet keys in the store.
const TimesheetKeyDateFmt = "2006-01-02"

// Timesheet represents a timesheet with entries.
type Timesheet struct {
	// The date of the timesheet.
	Key string
	// An array of the hashes of entries that belong in this timesheet.
	Entries []string
}

// NewTimesheet create a new instance of Timesheet.
func NewTimesheet() Timesheet {
	return Timesheet{
		Key: time.Now().Format(TimesheetKeyDateFmt),
	}
}

// FromMessage reads a `proto.TrackingTimesheet` message into this Timesheet.
func (t *Timesheet) FromMessage(message *proto.TrackingTimesheet) {
	t.Key = message.Key
	t.Entries = message.Entries
}

// ToMessage converts this Timesheet to a `proto.TrackingTimesheet`.
func (t *Timesheet) ToMessage() *proto.TrackingTimesheet {
	return &proto.TrackingTimesheet{
		Key:     t.Key,
		Entries: t.Entries,
	}
}

// AppendEntry appends a reference to an entry to the timesheet.
func (t *Timesheet) AppendEntry(entry Entry) {
	t.Entries = append(t.Entries, entry.Hash)
}

// RemoveEntry removes a reference to an entry from the timesheet.
func (t *Timesheet) RemoveEntry(entry Entry) {
	index := -1

	for idx, hash := range t.Entries {
		if hash == entry.Hash {
			index = idx
			break
		}
	}

	if index >= 0 {
		// Remove the entry
		t.Entries = append(t.Entries[:index], t.Entries[index+1:]...)
	}
}
