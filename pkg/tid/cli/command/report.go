package command

import (
	"fmt"
	"text/template"
	"time"

	"github.com/SeerUK/tid/pkg/tid/cli/display"
	"github.com/SeerUK/tid/pkg/timeutil"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ReportCommand creates a command to view a timesheet report.
func ReportCommand(factory util.Factory, config types.TomlConfig) *console.Command {
	var date time.Time
	var end time.Time
	var format string
	var start time.Time
	var noSummary bool

	configure := func(def *console.Definition) {
		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&date),
			Spec:  "-d, --date=DATE",
			Desc:  "The exact date of a timesheet to show a report for.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&end),
			Spec:  "-e, --end=END",
			Desc:  "The end date of the report. (Default: today)",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&format),
			Spec:  "-f, --format=FORMAT",
			Desc:  "Output formatting string. Uses Go templates.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&start),
			Spec:  "-s, --start=START",
			Desc:  "The start date of the report. (Default: today)",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewBoolValue(&noSummary),
			Spec:  "--no-summary",
			Desc:  "Hide the summary?",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		gateway := factory.BuildTrackingGateway()

		hasDate := input.HasOption([]string{"d", "date"})
		hasEnd := input.HasOption([]string{"e", "end"})
		hasFormat := input.HasOption([]string{"f", "format"})
		hasStart := input.HasOption([]string{"s", "start"})

		// We need to get the current date, this is a little hacky, but we need it without any time
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

		if len(entries) == 0 {
			output.Println("report: No entries within the given time period")
			return nil
		}

		if !noSummary {
			output.Printf("Report for %s.\n\n", getDateRange(start, end))
			output.Printf("Total Duration: %s\n", getDurationForEntries(entries))
			output.Printf("Entry Count: %d\n", len(entries))
			output.Println()
		}

		if hasFormat {
			for _, entry := range entries {
				tmpl := template.Must(template.New("entry-list").Parse(format))
				tmpl.Execute(output.Writer, entry)

				output.Println()
			}

			return nil
		}

		display.WriteEntriesTable(entries, output.Writer, config)

		return nil
	}

	return &console.Command{
		Name:        "report",
		Alias:       "rep",
		Description: "Display a timesheet report.",
		Configure:   configure,
		Execute:     execute,
	}
}

func getDateRange(start time.Time, end time.Time) string {
	if start.Equal(end) {
		return end.Format(timeutil.DateFmt)
	}

	return fmt.Sprintf("%s to %s", start.Format(timeutil.DateFmt), end.Format(timeutil.DateFmt))
}

func getDurationForEntries(entries []types.Entry) time.Duration {
	var duration time.Duration

	for _, e := range entries {
		duration = duration + e.Duration
	}

	return duration
}
