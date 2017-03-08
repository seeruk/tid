package tracking

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

// KeyStatus is the key for the current tracking status in the store.
const KeyStatus = "status"

// SysGateway provides access to tid system data in the database.
type SysGateway interface {
	// FindOrCreateStatus attempts to find the current status, if one is not in the store then a new
	// status object is instantiated.
	FindOrCreateStatus() (types.Status, error)
	// PersistStatus persists a given status to the store.
	PersistStatus(status types.Status) error
}

// storeSysGateway is a functional SysGateway.
type storeSysGateway struct {
	// The underlying storage to access.
	store state.Store
}

// NewStoreSysGateway creates a new sys gateway.
func NewStoreSysGateway(store state.Store) SysGateway {
	return &storeSysGateway{
		store: store,
	}
}

func (g *storeSysGateway) FindOrCreateStatus() (types.Status, error) {
	status := types.NewStatus()
	message := &proto.SysStatus{}

	err := g.store.Read(KeyStatus, message)
	if err != nil && err != state.ErrStoreNilResult {
		return status, err
	}

	if err == nil {
		status.FromMessage(message)
	}

	return status, nil
}

func (g *storeSysGateway) PersistStatus(status types.Status) error {
	return g.store.Write(KeyStatus, status.ToMessage())
}
