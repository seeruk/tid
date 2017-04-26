package entry

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

// ListCommand creates a command to list timesheet entries.
func ListCommand(factory util.Factory) *console.Command {
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

		now := timeutil.Date(time.Now())

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

		if hasFormat {
			for _, entry := range entries {
				tmpl := template.Must(template.New("entry-list").Parse(format))
				tmpl.Execute(output.Writer, entry)

				output.Println()
			}

			return nil
		}

		if len(entries) == 0 {
			return errors.New("list: No entries within the given time period")
		}

		display.WriteEntriesTable(entries, output.Writer)

		return nil
	}

	return &console.Command{
		Name:        "list",
		Alias:       "ls",
		Description: "List timesheet entries.",
		Configure:   configure,
		Execute:     execute,
	}
}
