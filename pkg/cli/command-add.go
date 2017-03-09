package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// AddCommand creates a command to add timesheet entries.
func AddCommand(gateway tracking.TimesheetGateway) *console.Command {
	var startedAt time.Time
	var duration time.Duration
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewDateValue(&startedAt),
			"STARTED",
			"When did you start working?",
		)

		def.AddArgument(
			parameters.NewDurationValue(&duration),
			"DURATION",
			"How long did you spend on what you want to add?",
		)

		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What were you working on?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		sheet, err := gateway.FindOrCreateTimesheet(startedAt.Format(types.TimesheetKeyDateFmt))
		if err != nil {
			return err
		}

		entry := types.NewEntry()
		entry.Duration = duration
		entry.Note = note
		entry.Timesheet = sheet.Key

		sheet.AppendEntry(entry)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistEntry(entry))
		errs.Add(gateway.PersistTimesheet(sheet))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Added entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "add",
		Description: "Add a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
