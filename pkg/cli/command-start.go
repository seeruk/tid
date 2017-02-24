package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/timesheet"
	"github.com/SeerUK/tid/proto"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

func StartCommand(gateway timesheet.Gateway) console.Command {
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What are you tracking?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		// @todo: Refactor into some kind of a facade (how can we store the result as a value in
		// that case, we'll have errors as values there still...)

		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		if timesheet.IsActive(status) {
			output.Println("start: Stop an existing timer before starting a new one")
			return nil
		}

		// @todo: This should be closer to the end.
		output.Printf("Started tracking '%s'.\n", note)

		// @todo: FindCurrentTimeSheet():
		now := time.Now().Local()

		sheet, err := gateway.FindTimeSheet(now)
		if err != nil {
			return err
		}

		// @todo: Maybe a helper for creating and appending instead of this?
		sheet.Entries = append(sheet.Entries, createEntry(note))

		// @todo: Make helper for this:
		status.State = proto.Status_STARTED
		status.TimeSheetEntry = &proto.TimeSheetEntryRef{
			Date:  now.Format(timesheet.KeyTimesheetFmt),
			Index: int64(len(sheet.Entries) - 1),
		}

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistStatus(&status))
		errs.Add(gateway.PersistTimesheet(now, &sheet))

		return errs.Errors()
	}

	return console.Command{
		Name:        "start",
		Description: "Start a new timer.",
		Configure:   configure,
		Execute:     execute,
	}
}

// @todo: This should not be in here.
func createEntry(note string) *proto.TimeSheetEntry {
	now := time.Now().Unix()

	return &proto.TimeSheetEntry{
		Note:      note,
		StartTime: uint64(now),
	}
}
