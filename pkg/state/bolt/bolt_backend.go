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

func (b *boltBackend) ForEach(bucket string, fn func(key string, val []byte) error) error {
	// @todo: This is horrific, but works. We can't call `fn` in the update / bucket's ForEach
	// because we end up locking the database when reading and can't perform writes like we might
	// want to when we call this method. A solution similar to this might be necessary, but surely
	// we can be more efficient and still do it one at a time... right? The problem is more to do
	// with loading the values into memory than the keys in our case at least.

	var keys [][]byte
	var vals [][]byte

	err := b.db.View(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.ForEach(func(k, v []byte) error {
			keys = append(keys, k)
			vals = append(vals, v)

			return nil
		})
	})

	if err != nil {
		return err
	}

	for i, key := range keys {
		err := fn(string(key), vals[i])
		if err != nil {
			return err
		}
	}

	return nil
}
