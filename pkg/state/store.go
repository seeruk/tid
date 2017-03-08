package state

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var (
	// ErrNilValue is the error given when a value passed in is nil.
	ErrNilValue = errors.New("state: `value` must not be null")
	// ErrNilResult is the error given when there is no entry found for a key in the database.
	ErrNilResult = errors.New("state: No value found")
)

// Store provides a means of persisting some data in a key/value store.
type Store interface {
	// Read a value with a given key from the store into a given message.
	Read(key string, value proto.Message) error
	// Write a given value to a key given in the store.
	Write(key string, value proto.Message) error
	// Delete a value with a given key from the store.
	Delete(key string) error
}
