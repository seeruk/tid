package cli

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// WorkspaceCommand creates a command for listing and switching workspaces.
func WorkspaceCommand(tsGateway tracking.TimesheetGateway) console.Command {
	var workspaceName string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&workspaceName),
			"[NAME]",
			"An optional workspace to switch to.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		output.Println(workspaceName)

		return nil
	}

	return console.Command{
		Name:        "workspace",
		Description: "List or switch workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
