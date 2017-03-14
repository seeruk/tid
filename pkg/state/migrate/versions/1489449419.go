package versions

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/types"
)

func init() {
	migrate.RegisterMigration(&Migration1489449419{})
}

// Migration1489449419 is a backend migration created at 1489449419 unix time.
type Migration1489449419 struct{}

// Description provides a description of what the migration is doing.
func (m *Migration1489449419) Description() string {
	return "Set up basic data structure."
}

// Migrate performs the migration.
func (m *Migration1489449419) Migrate(backend state.Backend) error {
	errs := errhandling.NewErrorStack()
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketSys))
	errs.Add(backend.CreateBucketIfNotExists(state.BackendBucketTimesheet))
	errs.Add(backend.CreateBucketIfNotExists(fmt.Sprintf(
		state.BackendBucketWorkspaceFmt,
		types.TrackingStatusDefaultWorkspace,
	)))

	return errs.Errors()
}

// Version returns the version number of the migration.
func (m *Migration1489449419) Version() uint {
	return 1489449419
}
