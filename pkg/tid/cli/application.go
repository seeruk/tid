package cli

import (
	"fmt"

	"github.com/SeerUK/tid/pkg/tid/cli/command"
	"github.com/SeerUK/tid/pkg/tid/cli/command/entry"
	"github.com/SeerUK/tid/pkg/tid/cli/command/timesheet"
	"github.com/SeerUK/tid/pkg/tid/cli/command/workspace"
	"github.com/eidolon/console"
)

var (
	// BuildTime should be set to a datetime string.
	BuildTime = "n/a"
	// Commit should be set to a Git commit SHA.
	Commit = "n/a"
	// Version should be set to the tid version.
	Version = "n/a"
)

// CreateApplication builds the console application instance. Providing it with some basic
// information like the name and version.
func CreateApplication(kernel *TidKernel) *console.Application {
	version := fmt.Sprintf("%s (%s, %s)", Version, Commit, BuildTime)

	application := console.NewApplication("tid", version)
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
			entry.CreateCommand(kernel.Factory),
			entry.DeleteCommand(kernel.Factory),
			entry.ListCommand(kernel.Factory, kernel.Config),
			entry.UpdateCommand(kernel.Factory),
		}),

		// Timesheet commands
		timesheet.RootCommand().AddCommands([]*console.Command{
			timesheet.DeleteCommand(kernel.Factory),
			timesheet.ListCommand(kernel.Factory, kernel.Config),
		}),

		// Workspace commands
		workspace.RootCommand().AddCommands([]*console.Command{
			workspace.CreateCommand(kernel.Factory),
			workspace.DeleteCommand(kernel.Factory),
			workspace.ListCommand(kernel.Factory),
			workspace.SwitchCommand(kernel.Factory),
		}),

		command.ReportCommand(kernel.Factory, kernel.Config),
		command.ResumeCommand(kernel.Factory),
		command.StartCommand(kernel.Factory),
		command.StatusCommand(kernel.Factory, kernel.Config),
		command.StopCommand(kernel.Factory),
	}
}
