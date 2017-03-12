package cli

import (
	"github.com/SeerUK/tid/pkg/cli/command/entry"
	"github.com/eidolon/console"
)

// GetCommands gets all of the commands registered in the application.
func GetCommands(kernel *TidKernel) []*console.Command {
	// @todo: Shouldn't need these when refactoring is complete, we pass in the factory.
	trackingSysGateway := kernel.TrackingFactory.BuildSysGateway()
	trackingTimesheetGateway := kernel.TrackingFactory.BuildTimesheetGateway()

	return []*console.Command{
		// Entry commands
		entry.RootCommand().AddCommands([]*console.Command{
			entry.ListCommand(kernel.TrackingFactory), // @todo: Write this
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

		// Top-level tracking commands
		// @todo: These should come from the `command` package.
		ReportCommand(trackingSysGateway, trackingTimesheetGateway),
		ResumeCommand(trackingSysGateway, trackingTimesheetGateway),
		StartCommand(trackingSysGateway, trackingTimesheetGateway),
		StatusCommand(trackingSysGateway, trackingTimesheetGateway),
		StopCommand(trackingSysGateway, trackingTimesheetGateway),
	}
}
