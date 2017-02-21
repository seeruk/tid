package state

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

var ErrNilValue = errors.New("state: `value` must not be null.")

// BoltStore implements the Store interface to provide a simple, fast, and reliable key / value
// store using Bolt.
type BoltStore struct {
	bucket string
	db     *bolt.DB
}

// NewBoltStore creates a new BoltStore instance.
func NewBoltStore(db *bolt.DB, bucket string) *BoltStore {
	return &BoltStore{
		bucket: bucket,
		db:     db,
	}
}

// Read reads a value from the BoltStore, and returns a ProtoBuf message.
func (b *BoltStore) Read(key string, value proto.Message) error {
	if value == nil {
		return ErrNilValue
	}

	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		result := bucket.Get([]byte(key))

		proto.Unmarshal(result, value)

		return nil
	})
}

// Write writes a value to a given key in the BoltStore.
func (b *BoltStore) Write(key string, value proto.Message) error {
	if value == nil {
		return ErrNilValue
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		result, err := proto.Marshal(value)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), result)
		if err != nil {
			return err
		}

		return nil
	})
}
