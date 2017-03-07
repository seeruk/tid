package cli

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

// StopCommand creates a command to stop timers.
func StopCommand(sysGateway tracking.SysGateway, tsGateway tracking.TimesheetGateway) console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if !status.IsRunning {
			output.Println("stop: There is no active timer running")
			return nil
		}

		entry, err := tsGateway.FindEntry(status.Entry)
		if err != nil {
			return err
		}

		entry.UpdateDuration()

		status.Stop()

		errs := errhandling.NewErrorStack()
		errs.Add(sysGateway.PersistStatus(status))
		errs.Add(tsGateway.PersistEntry(entry))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Stopped timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return console.Command{
		Name:        "stop",
		Description: "Stop an existing timer.",
		Execute:     execute,
	}
}
