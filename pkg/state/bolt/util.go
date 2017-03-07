package bolt

import (
	"fmt"
	"os"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/boltdb/bolt"
)

const (
	// BoltBucketSys is the bucket name for system information.
	BoltBucketSys = "tid_sys"
	// BoltBucketTimesheet is the original name of the Bucket for the timesheet stored in Bolt. This
	// is only used for migrating old data now.
	BoltBucketTimesheet = "tid_tracking"
	// BoltBucketWorkspaceDefault is the default workspace's name.
	BoltBucketWorkspaceDefault = "default"
	// BoltBucketWorkspaceFmt is the formatting string for timesheet bucket names.
	BoltBucketWorkspaceFmt = "tid_tracking_%s"
	// BoltDatabaseFilename is the name of the database file name on disk.
	BoltDatabaseFilename = "tid.db"
)

// Open opens a Bolt database, creating it if it doesn't exist already.
func Open(tidDir string) (*bolt.DB, error) {
	// Make the `path` if it does not exist.
	err := os.MkdirAll(tidDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return bolt.Open(fmt.Sprintf("%s/%s", tidDir, BoltDatabaseFilename), 0600, nil)
}

// Initialise ensures necessary Buckets are created and ready to use.
func Initialise(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists([]byte(BoltBucketSys))
		_, err2 := tx.CreateBucketIfNotExists([]byte(BoltBucketTimesheet))
		_, err3 := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf(
			BoltBucketWorkspaceFmt,
			BoltBucketWorkspaceDefault,
		)))

		errs := errhandling.NewErrorStack()
		errs.Add(err1)
		errs.Add(err2)
		errs.Add(err3)

		return errs.Errors()
	})

	if err != nil {
		return err
	}

	return migrateTimesheetToDefaultWorkspace(db)
}

// migrateTimesheetToDefaultWorkspace takes data in the old timesheet bucket and puts it in the new
// default workspace bucket.
func migrateTimesheetToDefaultWorkspace(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		sysBucket := tx.Bucket([]byte(BoltBucketSys))
		timesheetBucket := tx.Bucket([]byte(BoltBucketTimesheet))
		workspaceBucket := tx.Bucket([]byte(fmt.Sprintf(
			BoltBucketWorkspaceFmt,
			BoltBucketWorkspaceDefault,
		)))

		cursor := timesheetBucket.Cursor()
		errs := errhandling.NewErrorStack()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if string(k) == tracking.KeyStatus {
				errs.Add(sysBucket.Put(k, v))
			} else {
				errs.Add(workspaceBucket.Put(k, v))
			}
		}

		errs.Add(tx.DeleteBucket([]byte(BoltBucketTimesheet)))

		return errs.Errors()
	})
}
