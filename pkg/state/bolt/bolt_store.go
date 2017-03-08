package bolt

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/golang/protobuf/proto"

	boltdb "github.com/boltdb/bolt"
)

// boltStore implements the Store interface to provide a simple, fast, and reliable key / value
// store using Bolt.
type boltStore struct {
	bucket string
	db     *boltdb.DB
}

// NewBoltStore creates a new Store instance using Bolt.
func NewBoltStore(db *boltdb.DB, bucket string) state.Store {
	return &boltStore{
		bucket: bucket,
		db:     db,
	}
}

func (b *boltStore) Read(key string, value proto.Message) error {
	if value == nil {
		return state.ErrNilValue
	}

	return b.db.View(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))
		result := bucket.Get([]byte(key))

		if result == nil {
			return state.ErrNilResult
		}

		return proto.Unmarshal(result, value)
	})
}

func (b *boltStore) Write(key string, value proto.Message) error {
	if value == nil {
		return state.ErrNilValue
	}

	return b.db.Update(func(tx *boltdb.Tx) error {
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

func (b *boltStore) Delete(key string) error {
	return b.db.Update(func(tx *boltdb.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		return bucket.Delete([]byte(key))
	})
}
