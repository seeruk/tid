package cli

import (
	"github.com/SeerUK/tid/pkg/tid/cli/command"
	"github.com/SeerUK/tid/pkg/tid/cli/command/entry"
	"github.com/SeerUK/tid/pkg/tid/cli/command/timesheet"
	"github.com/SeerUK/tid/pkg/tid/cli/command/workspace"
	"github.com/eidolon/console"
)

// CreateApplication builds the console application instance. Providing it with some basic
// information like the name and version.
func CreateApplication(kernel *TidKernel) *console.Application {
	application := console.NewApplication("tid", "0.2.0-alpha.1")
	application.Logo = `
######## ### #######
   ###   ###       ##
   ###   ###  ###  ##
   ###   ###  ###  ##
   ###   ###  ######
`

	application.AddCommands(buildCommands(kernel))

	return application
}

// buildCommands instantiates all of the commands registered in the application.
func buildCommands(kernel *TidKernel) []*console.Command {
	return []*console.Command{
		// Entry commands
		entry.RootCommand().AddCommands([]*console.Command{
			entry.ListCommand(kernel.Factory),
			entry.CreateCommand(kernel.Factory),
			entry.UpdateCommand(kernel.Factory),
			entry.DeleteCommand(kernel.Factory),
		}),

		// Timesheet commands
		timesheet.RootCommand().AddCommands([]*console.Command{
			timesheet.ListCommand(kernel.Factory),
			timesheet.DeleteCommand(kernel.Factory),
		}),

		// Workspace commands
		// @todo: Write these
		workspace.RootCommand().AddCommands([]*console.Command{
			workspace.ListCommand(kernel.Factory),
		}),

		command.ReportCommand(kernel.Factory),
		command.ResumeCommand(kernel.Factory),
		command.StartCommand(kernel.Factory),
		command.StatusCommand(kernel.Factory),
		command.StopCommand(kernel.Factory),
	}
}
