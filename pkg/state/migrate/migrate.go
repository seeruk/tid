package migrate

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

// versions references all registered migrations.
var migrations []Migration

// Migration is the interface all backend migrations will follow.
type Migration interface {
	// Description provides a description of what the migration is doing.
	Description() string
	// Migrate performs the migration.
	Migrate(backend state.Backend) error
	// Version returns the version number of the migration.
	Version() uint
}

// RegisterMigration registers the given migration with the application.
func RegisterMigration(migration Migration) {
	migrations = append(migrations, migration)
}

// Backend ensures that the state.Backend is ready to use, and up-to-date.
func Backend(backend state.Backend) error {
	message := proto.SysMigrationsStatus{}

	store := state.NewBackendStore(backend, state.BackendBucketSys)
	store.Read(state.KeyMigrations, &message)

	status := types.NewMigrationsStatus()
	status.FromMessage(&message)

	//fmt.Println(status.LatestVersion())
	//
	//for _, migration := range migrations {
	//	fmt.Println(migration.Description())
	//}

	errs := errhandling.NewErrorStack()
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketSys))
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketTimesheet))
	errs.Add(backend.CreateBucketIfNotExists(fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		types.TrackingStatusDefaultWorkspace,
	)))

	if !errs.Empty() {
		return errs.Errors()
	}

	return migrateMonoBucketTimesheet(backend)
}

// migrateMonoBucketTimesheet takes data in the old timesheet bucket and puts it in the new
// default workspace bucket.
func migrateMonoBucketTimesheet(backend state.Backend) error {
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
