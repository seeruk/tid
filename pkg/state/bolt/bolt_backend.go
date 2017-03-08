package bolt

import (
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
	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.Put([]byte(key), value)
	})
}

func (b *boltBackend) Delete(bucket string, key string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.Delete([]byte(key))
	})
}

func (b *boltBackend) ForEach(bucket string, fn func(key string, val []byte) error) error {
	return b.db.View(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			err := fn(string(k), v)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
