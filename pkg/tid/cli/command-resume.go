package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ResumeCommand creates a command to resume timers.
func ResumeCommand(sysGateway tracking.SysGateway, tsGateway tracking.TimesheetGateway) *console.Command {
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&hash),
			"[HASH]",
			"A short or long hash for an entry.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.IsRunning {
			output.Println("resume: Stop an existing timer before resuming a one")
			return nil
		}

		if hash == "" {
			if status.Entry == "" {
				output.Println("resume: No timer to resume")
				return nil
			}

			hash = status.Entry
		}

		entry, err := tsGateway.FindEntry(hash)
		if err != nil {
			return err
		}

		sheet, err := tsGateway.FindOrCreateTimesheet(entry.Timesheet)
		if err != nil {
			return err
		}

		// Update the time that this entry was last updated.
		entry.Updated = time.Now()

		status.Start(sheet, entry)

		errs := errhandling.NewErrorStack()
		errs.Add(sysGateway.PersistStatus(status))
		errs.Add(tsGateway.PersistEntry(entry))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Resumed timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "resume",
		Description: "Resume an existing timer.",
		Configure:   configure,
		Execute:     execute,
	}
}