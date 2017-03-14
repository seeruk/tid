package versions

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/types"
)

func init() {
	migrate.RegisterMigration(&Migration1489498859{})
}

// Migration1489498859 is a backend migration created at 1489498859 unix time.
type Migration1489498859 struct{}

// Description provides a description of what the migration is doing.
func (m *Migration1489498859) Description() string {
	return "Migrate data into new bucket layout for workspaces."
}

// Migrate performs the migration.
func (m *Migration1489498859) Migrate(backend state.Backend) error {
	workspaceBucketName := fmt.Sprintf(state.BackendBucketWorkspaceFmt, types.TrackingStatusDefaultWorkspace)

	errs := errhandling.NewErrorStack()

	err := backend.ForEachSingle(state.BackendBucketTimesheet, func(key string, val []byte) error {
		if key == state.KeyStatus {
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

// Version returns the version number of the migration.
func (m *Migration1489498859) Version() uint {
	return 1489498859
}
