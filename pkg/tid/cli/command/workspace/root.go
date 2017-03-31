package workspace

import "github.com/eidolon/console"

// RootCommand creates a new command that is the parent to all workspace-related sub-commands.
func RootCommand() *console.Command {
	return &console.Command{
		Name:        "workspace",
		Alias:       "w",
		Description: "Manage workspaces.",
	}
}
