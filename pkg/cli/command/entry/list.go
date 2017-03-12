package entry

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

// ListCommand creates a command to list timesheet entries.
func ListCommand(factory tracking.Factory) *console.Command {
	return &console.Command{
		Name:        "list",
		Description: "List timesheet entries.",
	}
}
