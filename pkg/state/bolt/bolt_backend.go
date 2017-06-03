package bolt

import (
	"bytes"

	"github.com/SeerUK/tid/pkg/state"

	boltdb "github.com/boltdb/bolt"
)

// boltBackend implements the Backend interface to provide a simple, fast, and reliable key / value
// store, embedded within tid, using Bolt DB.
type boltBackend struct {
	db *boltdb.DB
}

// NewBoltBackend create a new Backend instance using Bolt.
func NewBoltBackend(db *boltdb.DB) state.Backend {
	return &boltBackend{
		db: db,
	}
}

func (b *boltBackend) CreateBucketIfNotExists(name string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))

		return err
	})
}

func (b *boltBackend) HasBucket(name string) bool {
	result := false

	b.db.View(func(tx *boltdb.Tx) error {
		result = tx.Bucket([]byte(name)) != nil

		return nil
	})

	return result
}

func (b *boltBackend) DeleteBucket(name string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		err := tx.DeleteBucket([]byte(name))

		if err == boltdb.ErrBucketNotFound {
			return state.ErrNilBucket
		}

		return err
	})
}

func (b *boltBackend) Read(bucket string, key string) ([]byte, error) {
	var value []byte

	err := b.db.View(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		if bucket == nil {
			return state.ErrNilBucket
		}

		value = bucket.Get([]byte(key))

		return nil
	})

	if err != nil {
		return value, err
	}

	if value == nil {
		return value, state.ErrStoreNilResult
	}

	return value, nil
}

func (b *boltBackend) Write(bucket string, key string, value []byte) error {
	err := b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.Put([]byte(key), value)
	})

	return err
}

func (b *boltBackend) Delete(bucket string, key string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.Delete([]byte(key))
	})
}

func (b *boltBackend) ForEachSingle(bucket string, fn func(key string, val []byte) error) error {
	// Load each key individually, get the value for the key, move onto next if there is any. This
	// should held ensure databases with large sets of keys are handled without crashing.

	var next []byte
	var last []byte
	var value []byte

	for len(next) == 0 || !bytes.Equal(next, last) {
		// Get next values
		err := b.db.View(func(tx *boltdb.Tx) error {
			bucket := tx.Bucket([]byte(bucket))
			cursor := bucket.Cursor()

			if len(next) == 0 {
				next, value = cursor.First()
				last, _ = cursor.Last()
			} else {
				cursor.Seek(next)

				next, value = cursor.Next()
			}

			return nil
		})

		if len(next) == 0 || len(last) == 0 {
			return nil
		}

		if err != nil {
			return err
		}

		// Run the user-defined function the item we're iterating over
		err = fn(string(next), value)
		if err != nil {
			return err
		}
	}

	return nil
}
