package cli

import (
	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StartCommand creates a command to start timers.
func StartCommand(sysGateway state.SysGateway, tsGateway state.TimesheetGateway) *console.Command {
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What are you working on?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.IsRunning {
			output.Println("start: Stop your existing timer before starting a new one")
			return nil
		}

		sheet, err := tsGateway.FindOrCreateTodaysTimesheet()
		if err != nil {
			return err
		}

		entry := types.NewEntry()
		entry.Note = note
		entry.Timesheet = sheet.Key

		sheet.AppendEntry(entry)

		status.Start(sheet, entry)

		errs := errhandling.NewErrorStack()
		errs.Add(sysGateway.PersistStatus(status))
		errs.Add(tsGateway.PersistEntry(entry))
		errs.Add(tsGateway.PersistTimesheet(sheet))

		if err = errs.Errors(); err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Started timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "start",
		Description: "Start a new timer.",
		Configure:   configure,
		Execute:     execute,
	}
}
