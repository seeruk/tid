package entry

import (
	"html/template"
	"time"

	"github.com/SeerUK/tid/pkg/cli/display"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ListDateFmt is the date format for report date ranges.
const ListDateFmt = "2006-01-02"

// ListCommand creates a command to list timesheet entries.
func ListCommand(factory tracking.Factory) *console.Command {
	var date time.Time
	var end time.Time
	var format string
	var start time.Time

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewDateValue(&date),
			"-d, --date=DATE",
			"The exact date of a timesheet to show a report for.",
		)

		def.AddOption(
			parameters.NewDateValue(&end),
			"-e, --end=END",
			"The end date of the report.",
		)

		def.AddOption(
			parameters.NewStringValue(&format),
			"-f, --format=FORMAT",
			"Output formatting string. Uses Go templates.",
		)

		def.AddOption(
			parameters.NewDateValue(&start),
			"-s, --start=START",
			"The start date of the report.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		gateway := factory.BuildTimesheetGateway()

		hasDate := input.HasOption([]string{"d", "date"})
		hasEnd := input.HasOption([]string{"e", "end"})
		hasFormat := input.HasOption([]string{"f", "format"})
		hasStart := input.HasOption([]string{"s", "start"})

		// We need to get the current date, this is a little hacky, but we need it without any time
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

		if hasDate {
			start = date
			end = date
		}

		entries, err := gateway.FindEntriesInDateRange(start, end)
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			output.Println("list: No entries within the given time period")
			return nil
		}

		if hasFormat {
			for _, entry := range entries {
				tmpl := template.Must(template.New("entry-list").Parse(format))
				tmpl.Execute(output.Writer, entry)

				output.Println()
			}

			return nil
		}

		display.WriteTableOfEntries(entries, output.Writer)

		return nil
	}

	return &console.Command{
		Name:        "list",
		Description: "List timesheet entries.",
		Configure:   configure,
		Execute:     execute,
	}
}
