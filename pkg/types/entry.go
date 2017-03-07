package types

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/SeerUK/tid/proto"
)

var (
	timeFormatShort = "3:04:05PM"
	timeFormatLong  = "3:04:05PM (2006-01-02)"
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
	// Whether or not this entry's timer is running.
	IsRunning bool
}

// NewEntry creates a new instance of Entry, with a new random hash, and dates set.
func NewEntry() Entry {
	return Entry{
		Hash:    createHash(),
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

// CreatedTimeFormat returns an appropriate time format for reporting output that is longer if the
// entry's created date was not the same as the timesheet it belongs to's date.
func (e Entry) CreatedTimeFormat() string {
	if e.Created.Format("2006-01-02") == e.Timesheet {
		return timeFormatShort
	}

	return timeFormatLong
}

// UpdatedTimeFormat returns an appropriate time format for reporting output that is longer if the
// entry's updated date was not the same as the timesheet it belongs to's date.
func (e Entry) UpdatedTimeFormat() string {
	if e.Updated.Format("2006-01-02") == e.Timesheet {
		return timeFormatShort
	}

	return timeFormatLong
}

// createHash creates a new random SHA-1 hash.
func createHash() string {
	nowUnix := time.Now().UnixNano()
	number := rand.Int()
	pid := os.Getpid()

	data := fmt.Sprintf("%d%d%d", nowUnix, number, pid)
	hash := sha1.Sum([]byte(data))

	return fmt.Sprintf("%x", hash)
}
