package state

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

const BoltBucketTimeSheet = "tid_time_sheet"
const BoltDatabaseFilename = "bolt.db"

// OpenBolt opens a Bolt database, creating it if it doesn't exist already.
func OpenBolt(tidDir string) (*bolt.DB, error) {
	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return bolt.Open(fmt.Sprintf("%s/%s", tidDir, BoltDatabaseFilename), 0600, nil)

}

// InitialiseBolt ensures necessary Buckets are created and ready to use.
func InitialiseBolt(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BoltBucketTimeSheet))

		return err
	})
}
