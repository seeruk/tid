#!/usr/bin/env bash

SCRIPT_DIR="$(dirname "$0")"
MIGRATIONS_DIR="$SCRIPT_DIR/../pkg/state/migrate/versions"

TIMESTAMP=$(date +%s)

CONTENT=$(cat <<MIG
package versions

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/state/migrate"
)

func init() {
	migrate.RegisterMigration(&Migration${TIMESTAMP}{})
}

// Migration${TIMESTAMP} is a backend migration created at ${TIMESTAMP} unix time.
type Migration${TIMESTAMP} struct{}

// Description provides a description of what the migration is doing.
func (m *Migration${TIMESTAMP}) Description() string {
	return "An incomplete migration."
}

// Migrate performs the migration.
func (m *Migration${TIMESTAMP}) Migrate(backend state.Backend) error {
	return nil
}

// Version returns the version number of the migration.
func (m *Migration${TIMESTAMP}) Version() uint {
	return ${TIMESTAMP}
}
MIG
)

echo "$CONTENT" > "$MIGRATIONS_DIR/$TIMESTAMP.go"
