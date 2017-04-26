package timesheet

import (
	"errors"
	"text/template"
	"time"

	"github.com/SeerUK/tid/pkg/tid/cli/display"
	"github.com/SeerUK/tid/pkg/timeutil"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ListCommand creates a command to list timesheets.
func ListCommand(factory util.Factory) *console.Command {
	var end time.Time
	var format string
	var start time.Time

	configure := func(def *console.Definition) {
		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&end),
			Spec:  "-e, --end=END",
			Desc:  "The end date of the listing. (Default: today)",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&format),
			Spec:  "-f, --format=FORMAT",
			Desc:  "Output formatting string. Uses Go templates.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&start),
			Spec:  "-s, --start=START",
			Desc:  "The start date of the listing. (Default: last monday)",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		trGateway := factory.BuildTrackingGateway()

		hasEnd := input.HasOption([]string{"e", "end"})
		hasFormat := input.HasOption([]string{"f", "format"})
		hasStart := input.HasOption([]string{"s", "start"})

		now := timeutil.Date(time.Now())

		if !hasStart {
			start = timeutil.LastWeekday(time.Monday)
		}

		if !hasEnd {
			end = now
		}

		// Lets be flexible... we'll get all by default, or we can use a range to limit output.
		ts, err := trGateway.FindTimesheetsInDateRange(start, end)
		if err != nil {
			return err
		}

		if hasFormat {
			for _, t := range ts {
				tmpl := template.Must(template.New("entry-list").Parse(format))
				tmpl.Execute(output.Writer, t)

				output.Println()
			}

			return nil
		}

		if len(ts) == 0 {
			return errors.New("list: No timesheets within the given time period")
		}

		display.WriteTimesheetsTable(ts, output.Writer)

		return nil
	}

	return &console.Command{
		Name:        "list",
		Alias:       "ls",
		Description: "List timesheets.",
		Configure:   configure,
		Execute:     execute,
	}
}
