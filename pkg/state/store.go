package state

import (
	"github.com/golang/protobuf/proto"
)

// Store provides a means of persisting some data in a key/value store.
type Store interface {
	// Read a value with a given key from the store into a given message.
	Read(key string, value proto.Message) error
	// Write a given value to a key given in the store.
	Write(key string, value proto.Message) error
}
