package cli

import (
	"github.com/SeerUK/tid/pkg/tid/cli/command"
	"github.com/SeerUK/tid/pkg/tid/cli/command/entry"
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
	// @todo: Shouldn't need these when refactoring is complete, we pass in the factory.
	trackingSysGateway := kernel.TrackingFactory.BuildSysGateway()
	trackingTimesheetGateway := kernel.TrackingFactory.BuildTimesheetGateway()

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
		// timesheet.RootCommand().AddCommands([]*console.Command{
		//
		// }),

		// Workspace commands
		// @todo: Write these
		// workspace.RootCommand().AddCommands([]*console.Command{
		//
		// }),

		command.ReportCommand(kernel.TrackingFactory),
		command.StatusCommand(kernel.TrackingFactory),

		// Top-level tracking commands
		// @todo: These should come from the `command` package.
		ResumeCommand(trackingSysGateway, trackingTimesheetGateway),
		StartCommand(trackingSysGateway, trackingTimesheetGateway),
		StopCommand(trackingSysGateway, trackingTimesheetGateway),
	}
}
