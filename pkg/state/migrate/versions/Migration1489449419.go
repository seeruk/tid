package versions

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
)

func init() {
	migrate.RegisterMigration(&Migration1489449419{})
}

// Migration1489449419 is a backend migration created at 1489449419 unix time.
type Migration1489449419 struct{}

// Description provides a description of what the migration is doing.
func (m *Migration1489449419) Description() string {
	return "A basic migration."
}

// Migrate performs the migration.
func (m *Migration1489449419) Migrate(backend state.Backend) error {
	return nil
}

// Version returns the version number of the migration.
func (m *Migration1489449419) Version() uint {
	return 1489449419
}
