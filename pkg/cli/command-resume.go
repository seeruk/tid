package cli

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

// ResumeCommand creates a command to resume timers.
func ResumeCommand(gateway tracking.Gateway) console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.IsActive() {
			output.Println("resume: Stop an existing timer before resuming a one")
			return nil
		}

		if status.Ref() == nil || status.Ref().Entry == "" {
			output.Println("resume: No timer to resume")
			return nil
		}

		sheet, err := gateway.FindOrCreateTimesheet(status.Ref().Timesheet)
		if err != nil {
			return err
		}

		entry, err := gateway.FindEntry(status.Ref().Entry)
		if err != nil {
			return err
		}

		// Update the time that this entry was last updated.
		entry.Update()

		status.Start(sheet, entry)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistEntry(entry))
		errs.Add(gateway.PersistStatus(status))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Resumed tracking '%s' (%s)\n", entry.Note(), entry.ShortHash())

		return nil
	}

	return console.Command{
		Name:        "resume",
		Description: "Resume an existing timer.",
		Execute:     execute,
	}
}