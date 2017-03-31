package entry

import "github.com/eidolon/console"

// RootCommand creates a new command that is the parent to all entry-related sub-commands.
func RootCommand() *console.Command {
	return &console.Command{
		Name:        "entry",
		Alias:       "e",
		Description: "Manage timesheet entries.",
	}
}
