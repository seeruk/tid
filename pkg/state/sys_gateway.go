package state

import (
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/proto"
)

const (
	// KeyMigrations is the key for the applied migration versions in the store.
	KeyMigrations = "migration_versions"
	// KeyStatus is the key for the current tracking status in the store.
	KeyStatus = "status"
)

// SysGateway provides access to tid system data in the database.
type SysGateway interface {
	// FindOrCreateMigrationsStatus attempts to find the current migrations information, if it can't
	// find any in the store then a new types.Migrations object is instantiated.
	FindOrCreateMigrationsStatus() (types.MigrationsStatus, error)
	// FindOrCreateStatus attempts to find the current status, if one is not in the store then a new
	// types.Status object is instantiated.
	FindOrCreateStatus() (types.TrackingStatus, error)
	// PersistMigrations persists a given types.Migrations to the store.
	PersistMigrations(migrations types.MigrationsStatus) error
	// PersistStatus persists a given types.Status to the store.
	PersistStatus(status types.TrackingStatus) error
}

// storeSysGateway is a functional SysGateway.
type storeSysGateway struct {
	// The underlying storage to access.
	store Store
}

// NewStoreSysGateway creates a new SysGateway.
func NewStoreSysGateway(store Store) SysGateway {
	return &storeSysGateway{
		store: store,
	}
}

func (g *storeSysGateway) FindOrCreateMigrationsStatus() (types.MigrationsStatus, error) {
	migrations := types.NewMigrationsStatus()
	message := &proto.SysMigrationsStatus{}

	err := g.store.Read(KeyMigrations, message)
	if err != nil && err != ErrStoreNilResult {
		return migrations, err
	}

	if err == nil {
		migrations.FromMessage(message)
	}

	return migrations, nil
}

func (g *storeSysGateway) FindOrCreateStatus() (types.TrackingStatus, error) {
	status := types.NewTrackingStatus()
	message := &proto.SysTrackingStatus{}

	err := g.store.Read(KeyStatus, message)
	if err != nil && err != ErrStoreNilResult {
		return status, err
	}

	if err == nil {
		status.FromMessage(message)
	}

	return status, nil
}

func (g *storeSysGateway) PersistMigrations(migrations types.MigrationsStatus) error {
	return g.store.Write(KeyMigrations, migrations.ToMessage())
}

func (g *storeSysGateway) PersistStatus(status types.TrackingStatus) error {
	return g.store.Write(KeyStatus, status.ToMessage())
}
