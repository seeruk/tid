package bolt

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/golang/protobuf/proto"

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

func (b *boltBackend) Read(bucket string, key string, value proto.Message) error {
	if value == nil {
		return state.ErrNilValue
	}

	return b.db.View(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		result := bucket.Get([]byte(key))

		if result == nil {
			return state.ErrNilResult
		}

		return proto.Unmarshal(result, value)
	})
}

func (b *boltBackend) Write(bucket string, key string, value proto.Message) error {
	if value == nil {
		return state.ErrNilValue
	}

	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

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

func (b *boltBackend) Delete(bucket string, key string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		return bucket.Delete([]byte(key))
	})
}

func (b *boltBackend) Close() error {
	return b.db.Close()
}
