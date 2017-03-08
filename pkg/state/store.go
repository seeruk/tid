package state

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

var (
	// ErrStoreNilValue is the error given when a value passed in is nil.
	ErrStoreNilMessage = errors.New("state: `message` must not be null")
	// ErrStoreNilResult is the error given when there is no entry found for a key in the database.
	ErrStoreNilResult = errors.New("state: No value found")
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

// backendStore is a functional Store.
type backendStore struct {
	backend Backend
	bucket  string
}

// NewBackendStore create a new Backend-based Store instance.
func NewBackendStore(backend Backend, bucket string) Store {
	return &backendStore{
		backend: backend,
		bucket:  bucket,
	}
}

func (b *backendStore) Read(key string, message proto.Message) error {
	if message == nil {
		return ErrStoreNilMessage
	}

	value, err := b.backend.Read(b.bucket, key)
	if err != nil {
		return err
	}

	return proto.Unmarshal(value, message)
}

func (b *backendStore) Write(key string, message proto.Message) error {
	if message == nil {
		return ErrStoreNilMessage
	}

	value, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	return b.backend.Write(b.bucket, key, value)
}

func (b *backendStore) Delete(key string) error {
	return b.backend.Delete(b.bucket, key)
}
