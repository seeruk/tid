package cli

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StartCommand creates a command to start timers.
func StartCommand(gateway tracking.Gateway) console.Command {
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What are you working on?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.IsActive() {
			output.Println("start: Stop your existing timer before starting a new one")
			return nil
		}

		sheet, err := gateway.FindOrCreateTodaysTimesheet()
		if err != nil {
			return err
		}

		entry := tracking.NewEntry()
		entry.Note = note
		entry.Timesheet = sheet.Key()

		sheet.AppendEntry(entry)

		status.Start(sheet, entry)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistEntry(entry))
		errs.Add(gateway.PersistTimesheet(sheet))
		errs.Add(gateway.PersistStatus(status))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Started tracking '%s' (%s)\n", entry.Note, entry.ShortHash)

		return nil
	}

	return console.Command{
		Name:        "start",
		Description: "Start a new timer.",
		Configure:   configure,
		Execute:     execute,
	}
}
