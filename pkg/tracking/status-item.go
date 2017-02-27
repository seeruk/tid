package tracking

import (
	"time"

	"github.com/SeerUK/tid/proto"
)

// StatusItem represents an item shown in status or report output, with all of the information
// necessary to display output.
type StatusItem struct {
	Timesheet string
	Hash      string
	ShortHash string
	Created   time.Time
	Updated   time.Time
	Note      string
	Duration  time.Duration
	Running   bool
}

// NewStatusItem creates a new instance of StatusItem.
func NewStatusItem(entry *Entry, running bool) StatusItem {
	return StatusItem{
		Timesheet: entry.Timesheet(),
		Hash:      entry.Hash(),
		ShortHash: entry.ShortHash(),
		Created:   entry.Created(),
		Updated:   entry.Updated(),
		Note:      entry.Note(),
		Duration:  entry.Duration(),
		Running:   running,
	}
}

// ToMessage converts this Entry
func (i *StatusItem) ToMessage() proto.TrackingEntry {
	return proto.TrackingEntry{
		Key:       i.Hash,
		Timesheet: i.Timesheet,
		Note:      i.Note,
		Created:   uint64(i.Created.Unix()),
		Updated:   uint64(i.Updated.Unix()),
		Duration:  uint64(i.Duration.Seconds()),
	}
}
