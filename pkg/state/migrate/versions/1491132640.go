package versions

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

func init() {
	migrate.RegisterMigration(&Migration1491132640{})
}

// Migration1491132640 is a backend migration created at 1491132640 unix time.
type Migration1491132640 struct{}

// Description provides a description of what the migration is doing.
func (m *Migration1491132640) Description() string {
	return "Set up SysWorkspaceIndex with default workspace."
}

// Migrate performs the migration.
func (m *Migration1491132640) Migrate(backend state.Backend) error {
	message := &proto.SysWorkspaceIndex{}
	message.Workspaces = append(message.Workspaces, types.TrackingStatusDefaultWorkspace)

	store := state.NewBackendStore(backend, state.BackendBucketSys)

	return store.Write(state.KeyWorkspaceIndex, message)

}

// Version returns the version number of the migration.
func (m *Migration1491132640) Version() uint {
	return 1491132640
}
