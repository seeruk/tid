package cli

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

// StopCommand creates a command to stop timers.
func StopCommand(gateway tracking.Gateway) console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if !status.IsActive {
			output.Println("stop: There is no active timer running")
			return nil
		}

		entry, err := gateway.FindEntry(status.Entry)
		if err != nil {
			return err
		}

		entry.UpdateDuration()

		status.Stop()

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistEntry(entry))
		errs.Add(gateway.PersistStatus(status))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Stopped tracking '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return console.Command{
		Name:        "stop",
		Description: "Stop an existing timer.",
		Execute:     execute,
	}
}
