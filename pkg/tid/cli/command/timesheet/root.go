package timesheet

import "github.com/eidolon/console"

// RootCommand creates a new command that is the parent to all timesheet-related sub-commands.
func RootCommand() *console.Command {
	return &console.Command{
		Name:        "timesheet",
		Alias:       "t",
		Description: "Manage timesheets.",
	}
}
