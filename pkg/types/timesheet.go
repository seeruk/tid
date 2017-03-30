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
	// An array of entries that belong in this timesheet.
	Entries []Entry
}

// NewTimesheet create a new instance of Timesheet.
func NewTimesheet() Timesheet {
	return Timesheet{
		Key: time.Now().Format(TimesheetKeyDateFmt),
	}
}

// FromMessageWithEntries reads a `proto.TrackingTimesheet` message into this Timesheet.
func (t *Timesheet) FromMessageWithEntries(message *proto.TrackingTimesheet, entries []Entry) {
	t.Key = message.Key
	t.Entries = entries
}

// ToMessage converts this Timesheet to a `proto.TrackingTimesheet`.
func (t *Timesheet) ToMessage() *proto.TrackingTimesheet {
	var entryKeys []string

	for _, entry := range t.Entries {
		entryKeys = append(entryKeys, entry.Hash)
	}

	return &proto.TrackingTimesheet{
		Key:     t.Key,
		Entries: entryKeys,
	}
}

// AppendEntry appends a reference to an entry to the timesheet.
func (t *Timesheet) AppendEntry(entry Entry) {
	t.Entries = append(t.Entries, entry)
}

// RemoveEntry removes a reference to an entry from the timesheet.
func (t *Timesheet) RemoveEntry(entry Entry) {
	index := -1

	for idx, tsEntry := range t.Entries {
		if tsEntry.Hash == entry.Hash {
			index = idx
			break
		}
	}

	if index >= 0 {
		// Remove the entry
		t.Entries = append(t.Entries[:index], t.Entries[index+1:]...)
	}
}
