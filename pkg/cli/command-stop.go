package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/proto"
	"github.com/eidolon/console"
)

func StopCommand(gateway tracking.Gateway) console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		// @todo: Refactor into some kind of a facade.

		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		if !status.IsActive() {
			output.Println("stop: There is no existing timer running")
			return nil
		}

		date, err := time.Parse(tracking.KeyTimesheetFmt, status.TimeSheetEntry.Date)
		if err != nil {
			return err
		}

		sheet, err := gateway.FindTimeSheet(date)
		if err != nil {
			return err
		}

		tracking.UpdateEntryDuration(&sheet, status.TimeSheetEntry.Index)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistStatus(tracking.NewStatus(&proto.Status{})))
		errs.Add(gateway.PersistTimesheet(date, &sheet))

		return errs.Errors()
	}

	return console.Command{
		Name:        "stop",
		Description: "Stop an existing timer.",
		Execute:     execute,
	}
}
