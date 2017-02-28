package tracking

import (
	"time"

	"github.com/SeerUK/tid/pkg/hash"
	"github.com/SeerUK/tid/proto"
)

// Entry represents a timesheet entry.
type Entry struct {
	// The date of the timesheet this entry belongs to.
	Timesheet string
	// The hash of this entry.
	Hash string
	// The time that this entry was created.
	Created time.Time
	// The time that this entry was last updated.
	Updated time.Time
	// The note on this entry.
	Note string
	// The amount of time logged against this entry.
	Duration time.Duration
}

// NewEntry creates a new instance of Entry, with a new random hash, and dates set.
func NewEntry() Entry {
	return Entry{
		Hash:    hash.CreateHash(),
		Created: time.Now(),
		Updated: time.Now(),
	}
}

// FromMessage reads a `proto.TrackingEntry` message into this Entry.
func (e *Entry) FromMessage(message *proto.TrackingEntry) {
	e.Timesheet = message.Timesheet
	e.Hash = message.Key
	e.Created = time.Unix(int64(message.Created), 0)
	e.Updated = time.Unix(int64(message.Updated), 0)
	e.Note = message.Note
	e.Duration = time.Duration(message.Duration) * time.Second
}

// ToMessage converts this Entry into a `proto.TrackingEntry`.
func (e *Entry) ToMessage() *proto.TrackingEntry {
	return &proto.TrackingEntry{
		Key:       e.Hash,
		Timesheet: e.Timesheet,
		Note:      e.Note,
		Created:   uint64(e.Created.Unix()),
		Updated:   uint64(e.Updated.Unix()),
		Duration:  uint64(e.Duration.Seconds()),
	}
}

// UpdateDuration adds the difference between the time this entry was last stopped and now to the
// duration. This also updates `Entry.Updated`.
func (e *Entry) UpdateDuration() {
	diff := time.Now().Sub(e.Updated)

	// We only care about the seconds, nothing more specific, otherwise output is too long.
	seconds := e.Duration.Seconds() + diff.Seconds()

	e.Duration = time.Duration(seconds) * time.Second
	e.Updated = time.Now()
}

// ShortHash returns a shortened version of this Entry's hash.
func (e Entry) ShortHash() string {
	return e.Hash[:7]
}
