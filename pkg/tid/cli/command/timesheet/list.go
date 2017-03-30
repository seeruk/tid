package timesheet

import (
	"text/template"
	"time"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ListDateFmt is the date format for report date ranges.
const ListDateFmt = "2006-01-02"

// ListCommand creates a command to list timesheets.
func ListCommand(factory tracking.Factory) *console.Command {
	var end time.Time
	var format string
	var start time.Time

	configure := func(def *console.Definition) {
		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&end),
			Spec:  "-e, --end=END",
			Desc:  "The end date of the report.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&format),
			Spec:  "-f, --format=FORMAT",
			Desc:  "Output formatting string. Uses Go templates.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&start),
			Spec:  "-s, --start=START",
			Desc:  "The start date of the report.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		trGateway := factory.BuildTrackingGateway()

		hasEnd := input.HasOption([]string{"e", "end"})
		hasFormat := input.HasOption([]string{"f", "format"})
		hasStart := input.HasOption([]string{"s", "start"})

		now, err := time.Parse(ListDateFmt, time.Now().Format(ListDateFmt))
		if err != nil {
			return err
		}

		if !hasStart {
			start = now
		}

		if !hasEnd {
			end = now
		}

		var ts []types.Timesheet

		// Lets be flexible... we'll get all by default, or we can use a range to limit output.
		if hasStart || hasEnd {
			ts, err = trGateway.FindTimesheetsInDateRange(start, end)
		} else {
			ts, err = trGateway.FindTimesheets()
		}

		if err != nil {
			return err
		}

		if len(ts) == 0 {
			output.Println("list: No timesheets within the given time period")
			return nil
		}

		if hasFormat {
			for _, t := range ts {
				tmpl := template.Must(template.New("entry-list").Parse(format))
				tmpl.Execute(output.Writer, t)

				output.Println()
			}

			return nil
		}

		output.Println(ts)

		return nil
	}

	return &console.Command{
		Name:        "list",
		Description: "List timesheets.",
		Configure:   configure,
		Execute:     execute,
	}
}
