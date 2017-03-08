package bolt

import (
	"fmt"
	"os"

	boltdb "github.com/boltdb/bolt"
)

// BoltDatabaseFilename is the name of the database file name on disk.
const BoltDatabaseFilename = "tid.db"

// Open opens a Bolt database, creating it if it doesn't exist already.
func Open(tidDir string) (*boltdb.DB, error) {
	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return boltdb.Open(fmt.Sprintf("%s/%s", tidDir, BoltDatabaseFilename), 0600, nil)
}
