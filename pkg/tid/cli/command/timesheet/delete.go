package timesheet

import (
	"time"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// DeleteCommand creates a command that is used to delete timesheets.
func DeleteCommand(factory tracking.Factory) *console.Command {
	var date time.Time

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewDateValue(&date),
			Spec:  "DATE",
			Desc:  "The date of the timesheet.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildTimesheetFacade()

		sheet, err := facade.Delete(date)
		if err != nil {
			return err
		}

		output.Printf("Deleted timesheet '%s'\n", sheet.Key)

		return nil
	}

	return &console.Command{
		Name:        "delete",
		Alias:       "d",
		Description: "Delete a timesheet, and it's entries.",
		Configure:   configure,
		Execute:     execute,
	}
}
