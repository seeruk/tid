package cli

import (
	"github.com/eidolon/console"
)

// GetCommands gets all of the commands registered in the application.
func GetCommands(kernel *TidKernel) []*console.Command {
	// @todo: Refactor to just pass in the kernel.
	return []*console.Command{
		AddCommand(kernel.TimesheetGateway),
		EditCommand(kernel.SysGateway, kernel.TimesheetGateway),
		RemoveCommand(kernel.TimesheetGateway, kernel.Facade),
		ReportCommand(kernel.SysGateway, kernel.TimesheetGateway),
		ResumeCommand(kernel.SysGateway, kernel.TimesheetGateway),
		StartCommand(kernel.SysGateway, kernel.TimesheetGateway),
		StatusCommand(kernel.SysGateway, kernel.TimesheetGateway),
		StopCommand(kernel.SysGateway, kernel.TimesheetGateway),
		WorkspaceCommand(kernel.Backend, kernel.SysGateway).AddCommands([]*console.Command{
			WorkspaceDeleteCommand(),
		}),
	}
}
