package migrate

import (
	"sort"

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

	// This may not always exist, and will cause a panic if it doesn't otherwise.
	backend.CreateBucketIfNotExists(state.BackendBucketSys)

	store := state.NewBackendStore(backend, state.BackendBucketSys)
	store.Read(state.KeyMigrations, &message)

	status := types.NewMigrationsStatus()
	status.FromMessage(&message)

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})

	currentVersion := status.CurrentVersion()
	missingVersions := migrations

	if currentVersion != 0 {
		index := indexOf(len(migrations), func(i int) bool {
			return migrations[i].Version() == currentVersion
		})

		missingVersions = migrations[index+1:]
	}

	for _, migration := range missingVersions {
		err := migration.Migrate(backend)
		if err != nil {
			return err
		}

		status.Versions = append(status.Versions, migration.Version())
	}

	store.Write(state.KeyMigrations, status.ToMessage())

	return nil
}

// indexOf uses the callback to find the index for some value.
func indexOf(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}

	return -1
}
