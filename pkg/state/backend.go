package state

import (
	"errors"
)

const (
	// BackendBucketSys is the bucket name for system information.
	BackendBucketSys = "tid_sys"
	// BackendBucketTimesheet is the original name of the Bucket for timesheet data. This is only
	// used for migrating old data now.
	BackendBucketTimesheet = "tid_tracking"
	// BackendBucketWorkspaceFmt is the formatting string for timesheet bucket names.
	BackendBucketWorkspaceFmt = "tid_tracking_%s"
)

// ErrNilBucket is the error given when there is no entry found for a key in the database.
var ErrNilBucket = errors.New("state: No bucket found")

// Backend provides an abstraction over underlying backend database technologies. It is separate to
// Store because Store is a more specialised interface with less functionality. Backend does not
// deal with ProtoBuf messages.
type Backend interface {
	// -- Buckets
	// CreateBucketIfNotExists attempts to create a bucket with a given name, if it doesn't exist.
	CreateBucketIfNotExists(name string) error
	// HasBucket returns true if this Backend has a bucket with the given name.
	HasBucket(name string) bool
	// DeleteBucket attempts to remove a Bucket. Returns ErrNilBucket if the bucket doesn't exist.
	DeleteBucket(name string) error

	// -- Key/Value Pairs
	// Read a value with a given key from a given bucket into a given message.
	Read(bucket string, key string) ([]byte, error)
	// Write a given value to a key given in the store.
	Write(bucket string, key string, val []byte) error
	// Delete a value with a given key from the store.
	Delete(bucket string, key string) error
	// ForEachSingle loops over each key/value pair individually in the given bucket. As buckets can
	// contain different data types we resort to using byte arrays for values.
	ForEachSingle(bucket string, fn func(key string, val []byte) error) error
}
