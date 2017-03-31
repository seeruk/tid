package cli

import (
	"github.com/SeerUK/tid/pkg/tid/cli/command"
	"github.com/SeerUK/tid/pkg/tid/cli/command/entry"
	"github.com/SeerUK/tid/pkg/tid/cli/command/timesheet"
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
			entry.ListCommand(kernel.TrackingFactory),
			entry.CreateCommand(kernel.TrackingFactory),
			entry.UpdateCommand(kernel.TrackingFactory),
			entry.DeleteCommand(kernel.TrackingFactory),
		}),

		// Timesheet commands
		// @todo: Write these
		timesheet.RootCommand().AddCommands([]*console.Command{
			timesheet.ListCommand(kernel.TrackingFactory),
		}),

		// Workspace commands
		// @todo: Write these
		// workspace.RootCommand().AddCommands([]*console.Command{
		//
		// }),

		command.ReportCommand(kernel.TrackingFactory),
		command.ResumeCommand(kernel.TrackingFactory),
		command.StartCommand(kernel.TrackingFactory),
		command.StatusCommand(kernel.TrackingFactory),
		command.StopCommand(kernel.TrackingFactory),
	}
}
