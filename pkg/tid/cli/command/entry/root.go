package entry

import "github.com/eidolon/console"

// RootCommand creates a new command that is the parent to all Entry-related sub-commands.
func RootCommand() *console.Command {
	return &console.Command{
		Name:        "entry",
		Description: "Manage timesheet entries.",
	}
}
