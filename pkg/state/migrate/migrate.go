package migrate

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/pkg/types"
)

// Backend ensures that the state.Backend is ready to use, and up-to-date.
func Backend(backend state.Backend) error {
	errs := errhandling.NewErrorStack()
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketSys))
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketTimesheet))
	errs.Add(backend.CreateBucketIfNotExists(fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		types.StatusDefaultWorkspace,
	)))

	if !errs.Empty() {
		return errs.Errors()
	}

	return migrateMonoBucketTimesheet(backend)
}

// migrateTimesheetToDefaultWorkspace takes data in the old timesheet bucket and puts it in the new
// default workspace bucket.
func migrateMonoBucketTimesheet(backend state.Backend) error {
	workspaceBucketName := fmt.Sprintf(state.BackendBucketWorkspaceFmt, types.StatusDefaultWorkspace)

	errs := errhandling.NewErrorStack()

	err := backend.ForEach(state.BackendBucketTimesheet, func(key string, val []byte) error {
		if key == tracking.KeyStatus {
			errs.Add(backend.Write(state.BackendBucketSys, key, val))
		} else {
			errs.Add(backend.Write(workspaceBucketName, key, val))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !errs.Empty() {
		return errs.Errors()
	}

	return backend.DeleteBucket(state.BackendBucketTimesheet)
}
