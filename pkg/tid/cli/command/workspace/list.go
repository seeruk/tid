package workspace

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
)

// ListCommand creates a command to list available workspaces.
func ListCommand(factory util.Factory) *console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		sysGateway := factory.BuildSysGateway()

		errs := errhandling.NewErrorStack()

		index, err1 := sysGateway.FindWorkspaceIndex()
		status, err2 := sysGateway.FindOrCreateStatus()

		errs.Add(err1)
		errs.Add(err2)

		if !errs.Empty() {
			return errs.Errors()
		}

		for _, workspace := range index.Workspaces {
			fmt := "%s\n"

			if workspace == status.Workspace {
				fmt = "%s *\n"
			}

			output.Printf(fmt, workspace)
		}

		return nil
	}

	return &console.Command{
		Name:        "list",
		Alias:       "ls",
		Description: "List available workspaces.",
		Execute:     execute,
	}
}
