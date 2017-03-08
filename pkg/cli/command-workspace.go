package cli

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// WorkspaceCommand creates a command for listing and switching workspaces.
func WorkspaceCommand(sysGateway tracking.SysGateway) console.Command {
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

		// Potential new command reference? Each entity type has relevant commands, and then there
		// are shorter named commands for easy access to certain things.
		//
		// What would the entry list command do? Could it be an easier way to provide entries for
		// completions without going through timesheet? It could use Bolt's looping over keys, via
		// some kind of abstraction.
		//
		// $ tid entry [list]
		// $ tid entry create <DURATION> <NOTE> --date=<DATE>
		// $ tid entry edit <HASH> --offset=<OFFSET>
		// $ tid entry delete <HASH>
		//
		// $ tid timesheet [list]
		// $ tid timesheet delete <DATE>
		//
		// $ tid workspace [list]
		// $ tid workspace create <NAME>
		// $ tid workspace switch <NAME>
		// $ tid workspace delete <NAME>
		//
		// $ tid start <NOTE>
		// $ tid resume [<HASH>]
		// $ tid status [<HASH>]
		// $ tid stop
		// $ tid report

		if isSwitching {
			// How do we know if the workspace exists? Should we validate that? Should you have to
			// create a workspace first? Subcommands sound like a dream now...
			output.Println(fmt.Sprintf("Switching to '%s'.", workspaceName))

			// "Switched to {new,existing} workspace '%s'."
			//
			// This implies that we do know which workspaces exist and which don't. Also, the Bucket
			// for the workspace needs to exist! There needs to be something abstracted from Bolt
			// that allows you to switch or create a bucket. We can use the same terminology...
			//
			// For knowing if workspaces exist, this can be part of `tid_sys`.
		} else {
			// @todo list in table and allow use of `--format` option.
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
