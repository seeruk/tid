package cli

import (
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// WorkspaceDeleteCommand creates a command for deleting workspaces.
func WorkspaceDeleteCommand() *console.Command {
	var name string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&name),
			"NAME",
			"A workspace to delete.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		output.Printf("Deleted workspace '%s'.", name)

		return nil
	}

	return &console.Command{
		Name:        "delete",
		Description: "Delete a given workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
