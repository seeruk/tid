package tracking

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/SeerUK/tid/proto"
)

// Entry wraps a ProtoBuf-generated proto.TrackingEntry message with helper methods. No state should
// be kept in this struct.
type Entry struct {
	Message *proto.TrackingEntry
}

// NewEntry creates a new instance of Entry.
func NewEntry(note string) *Entry {
	nowUnix := time.Now().Unix()

	created := uint64(nowUnix)
	updated := uint64(nowUnix)

	key := createNewKey()

	return &Entry{
		Message: &proto.TrackingEntry{
			Key:     key,
			Note:    note,
			Created: created,
			Updated: updated,
		},
	}
}

// Created gets the created time of the underlying message.
func (e *Entry) Created() time.Time {
	return time.Unix(int64(e.Message.Created), 0)
}

// Duration gets the duration as a time.Duration for nice formatting.
func (e *Entry) Duration() time.Duration {
	return time.Duration(e.Message.Duration) * time.Second
}

// Hash returns the key of the underlying message.
func (e *Entry) Hash() string {
	return e.Message.Key
}

// Note gets the note of the underlying message.
func (e *Entry) Note() string {
	return e.Message.Note
}

// ShortHash returns the short version of the key.
func (e *Entry) ShortHash() string {
	return e.Message.Key[0:7]
}

// Update updates the Updated timestamp in the underlying message.
func (e *Entry) Update() {
	e.Message.Updated = uint64(time.Now().Unix())
}

// Updated gets the updated time of the underlying message.
func (e *Entry) Updated() time.Time {
	return time.Unix(int64(e.Message.Updated), 0)
}

// UpdateDuration adds the difference between the time this entry was last stopped and now to the
// duration on the underlying message.
func (e *Entry) UpdateDuration() {
	e.Message.Duration = e.Message.Duration + (uint64(time.Now().Unix()) - e.Message.Updated)
}

// createNewKey creates a
func createNewKey() string {
	nowUnix := time.Now().UnixNano()
	number := rand.Int()
	pid := os.Getpid()

	data := fmt.Sprintf("%d%d%d", nowUnix, number, pid)
	hash := sha1.Sum([]byte(data))

	return fmt.Sprintf("%x", hash)
}
