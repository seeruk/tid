package entry

import (
	"text/template"
	"time"

	"github.com/SeerUK/tid/pkg/tid/cli/display"
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
		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&date),
			Spec:  "-d, --date=DATE",
			Desc:  "The exact date of a timesheet to show a listing for.",
		})

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
			Desc:  "The start date of the listing. (Default: today)",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		gateway := factory.BuildTrackingGateway()

		hasDate := input.HasOption([]string{"d", "date"})
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

		display.WriteEntriesTable(entries, output.Writer)

		return nil
	}

	return &console.Command{
		Name:        "list",
		Description: "List timesheet entries.",
		Configure:   configure,
		Execute:     execute,
	}
}
