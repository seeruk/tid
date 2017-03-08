package cli

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// WorkspaceCommand creates a command for listing and switching workspaces.
func WorkspaceCommand(backend state.Backend, sysGateway tracking.SysGateway) console.Command {
	var workspaceName string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&workspaceName),
			"[NAME]",
			"An optional workspace to switch to.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		isSwitching := workspaceName != ""

		if isSwitching {
			output.Println(fmt.Sprintf("Switching to '%s'.", workspaceName))
		} else {
			output.Println(fmt.Sprintf("Workspace: %s", status.Workspace))
		}

		return nil
	}

	return console.Command{
		Name:        "workspace",
		Description: "List or switch workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
