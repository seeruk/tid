package versions

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/pkg/util"
)

func init() {
	migrate.RegisterMigration(&Migration1496518709{})
}

// Migration1496518709 is a backend migration created at 1496518709 unix time.
type Migration1496518709 struct{}

// Description provides a description of what the migration is doing.
func (m *Migration1496518709) Description() string {
	return "Set up status to point to the default workspace if no workspace is set."
}

// Migrate performs the migration.
func (m *Migration1496518709) Migrate(backend state.Backend) error {
	factory := util.NewStandardFactory(backend)
	sysGateway := factory.BuildSysGateway()

	status, err := sysGateway.FindOrCreateStatus()

	if err != nil {
		return err
	}

	if status.Workspace != "" {
		return nil
	}

	status.Workspace = types.TrackingStatusDefaultWorkspace

	return sysGateway.PersistStatus(status)
}

// Version returns the version number of the migration.
func (m *Migration1496518709) Version() uint {
	return 1496518709
}
