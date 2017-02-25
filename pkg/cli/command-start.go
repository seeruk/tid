package cli

import (
	"time"

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
			"What are you tracking?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		if status.IsActive() {
			output.Println("start: Stop an existing timer before starting a new one")
			return nil
		}

		sheet, err := gateway.FindCurrentTimeSheet()
		if err != nil {
			return err
		}

		// entryRef = timesheet.AppendNewEntry(note)
		tracking.AppendNewEntry(&sheet, note)

		status.Start(&sheet)

		// date, err := time.Parse(tracking.KeyTimesheetFmt, entryRef.Date)
		date, err := time.Parse(tracking.KeyTimeSheetFmt, status.TimeSheetEntry().Date)
		if err != nil {
			return err
		}

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistStatus(status))
		errs.Add(gateway.PersistTimesheet(date, &sheet))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Started tracking '%s'.\n", note)

		return nil
	}

	return console.Command{
		Name:        "start",
		Description: "Start a new timer.",
		Configure:   configure,
		Execute:     execute,
	}
}
