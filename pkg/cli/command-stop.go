package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/timesheet"
	"github.com/SeerUK/tid/proto"
	"github.com/eidolon/console"
)

func StopCommand(gateway timesheet.Gateway) console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		// @todo: Refactor into some kind of a facade.

		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		// @todo: Some IsActive helper?
		if status.State != proto.Status_STARTED && status.State != proto.Status_PAUSED {
			output.Println("stop: There is no existing timer running.")
			return nil
		}

		date, err := time.Parse(timesheet.KeyTimesheetFmt, status.TimeSheetEntry.Date)
		if err != nil {
			return err
		}

		sheet, err := gateway.FindTimeSheet(date)
		if err != nil {
			return err
		}

		// @todo: Some helper for the entry update? Pass in sheet, and index.
		entryIdx := status.TimeSheetEntry.Index

		entry := sheet.Entries[entryIdx]
		entry.Duration = uint64(time.Now().Unix()) - entry.StartTime

		timesheet.ResetStatus(&status)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistStatus(&status))
		errs.Add(gateway.PersistTimesheet(date, &sheet))

		return errs.Errors()
	}

	return console.Command{
		Name:        "stop",
		Description: "Stop an existing timer.",
		Execute:     execute,
	}
}
