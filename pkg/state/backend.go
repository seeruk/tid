package state

import (
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
)

// ErrNilBucket is the error given when there is no entry found for a key in the database.
var ErrNilBucket = errors.New("state: No bucket found")

// Backend provides an abstraction over underlying backend database technologies. It is separate to
// Store because Store is a more specialised interface with less functionality.
type Backend interface {
	// -- Buckets
	// CreateBucketIfNotExists attempts to create a bucket with a given name, if it doens't exist.
	CreateBucketIfNotExists(name string) error
	// HasBucket returns true if this Backend has a bucket with the given name.
	HasBucket(name string) bool
	// DeleteBucket attempts to remove a Bucket.
	DeleteBucket(name string) error

	// -- Key/Value Pairs
	// Read a value with a given key from a given bucket into a given message.
	Read(bucket string, key string, value proto.Message) error
	// Write a given value to a key given in the store.
	Write(bucket string, key string, value proto.Message) error
	// Delete a value with a given key from the store.
	Delete(bucket string, key string) error

	// -- Misc
	// Some Backends will require closing.
	io.Closer
}
