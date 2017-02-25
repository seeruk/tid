package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

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
		// @todo: Refactor into some kind of a facade (how can we store the result as a value in
		// that case, we'll have errors as values there still...)

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

		tracking.AppendNewEntry(&sheet, note)
		tracking.UpdateStatusStartEntry(status, &sheet)

		// Where does the gateway fall into this refactor?
		// Bucket name should be tracking. Gateway is fine again then.

		// timesheet = tracking.NewTimeSheet()
		// entryRef = timesheet.AppendNewEntry(note)

		// status = tracking.NewStatus()
		// status.Start(entryRef)

		errs := errhandling.NewErrorStack()
		errs.Add(gateway.PersistStatus(status))
		errs.Add(gateway.PersistTimesheet(time.Now().Local(), &sheet))

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
