package types

import (
	"sort"

	"github.com/SeerUK/tid/proto"
)

// MigrationsStatus represents the current state of backend migrations.
type MigrationsStatus struct {
	// Versions is an array of the applied migration versions.
	Versions []uint
}

// NewMigrationsStatus creates a new instance of MigrationsStatus.
func NewMigrationsStatus() MigrationsStatus {
	return MigrationsStatus{}
}

// FromMessage reads a `proto.SysMigrationsStatus` message into this MigrationsStatus.
func (m *MigrationsStatus) FromMessage(message *proto.SysMigrationsStatus) {
	for _, v := range message.Versions {
		m.Versions = append(m.Versions, uint(v))
	}
}

// ToMessage converts this MigrationsStatus into a `proto.SysMigrationsStatus`.
func (m *MigrationsStatus) ToMessage() *proto.SysMigrationsStatus {
	message := proto.SysMigrationsStatus{}

	for _, v := range m.Versions {
		message.Versions = append(message.Versions, uint64(v))
	}

	return &message
}

// LatestVersion gets the latest applied migration version.
func (m *MigrationsStatus) LatestVersion() uint {
	versions := make([]uint, len(m.Versions))

	copy(versions, m.Versions)

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	if len(versions) > 0 {
		return versions[len(versions)-1]
	}

	return 0
}
