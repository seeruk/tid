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

	store := state.NewBackendStore(backend, state.BackendBucketSys)
	store.Read(state.KeyMigrations, &message)

	status := types.NewMigrationsStatus()
	status.FromMessage(&message)

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})

	for _, migration := range migrations {
		err := migration.Migrate(backend)
		if err != nil {
			return err
		}

		status.Versions = append(status.Versions, migration.Version())
	}

	return nil
}
