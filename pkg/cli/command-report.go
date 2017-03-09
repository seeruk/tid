package cli

import (
	"fmt"
	"text/template"
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
	"github.com/olekukonko/tablewriter"
)

// ReportDateFmt is the date format for report date ranges.
const ReportDateFmt = "2006-01-02"

// reportOutputItem represents the formattable source of an item in the report command output.
type reportOutputItem struct {
	Entry  types.Entry
	Status types.Status
}

// ReportCommand creates a command to view a timesheet report.
func ReportCommand(sysGateway tracking.SysGateway, tsGateway tracking.TimesheetGateway) *console.Command {
	var start time.Time
	var end time.Time
	var date time.Time
	var format string
	var noSummary bool

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewDateValue(&start),
			"-s, --start=START",
			"The start date of the report.",
		)

		def.AddOption(
			parameters.NewDateValue(&end),
			"-e, --end=END",
			"The end date of the report.",
		)

		def.AddOption(
			parameters.NewDateValue(&date),
			"-d, --date=DATE",
			"The exact date of a timesheet to show a report for.",
		)

		def.AddOption(
			parameters.NewStringValue(&format),
			"-f, --format=FORMAT",
			"Format string, uses Go templates.",
		)

		def.AddOption(
			parameters.NewBoolValue(&noSummary),
			"--no-summary",
			"Hide the summary?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		hasStart := input.HasOption([]string{"s", "start"})
		hasEnd := input.HasOption([]string{"e", "end"})
		hasDate := input.HasOption([]string{"d", "date"})

		// We need to get the current date, this is a little hacky, but we need it without any time
		now, err := time.Parse(ReportDateFmt, time.Now().Format(ReportDateFmt))
		if err != nil {
			return err
		}

		if hasDate {
			start = date
			end = date
		} else {
			if !hasStart {
				start = now
			}

			if !hasEnd {
				end = now
			}
		}

		if start.After(end) {
			output.Println("report: The start date must be before the end date")
			return nil
		}

		keys := getDateRangeTimesheetKeys(start, end)
		sheets, err := getTimesheetsByKeys(tsGateway, keys)
		if err != nil {
			return err
		}

		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		var duration time.Duration
		var entries int

		err = forEachEntry(tsGateway, sheets, func(entry types.Entry) {
			if entry.IsRunning {
				entry.UpdateDuration()
				tsGateway.PersistEntry(entry)
			}

			duration = duration + entry.Duration
			entries = entries + 1
		})

		if err != nil {
			return err
		}

		if entries == 0 {
			output.Println("report: No entries within the given time period")
			return nil
		}

		if !noSummary {
			if start.Equal(end) {
				format := "Report for %s.\n"
				output.Printf(format, end.Format(ReportDateFmt))
				output.Println()
			} else {
				format := "Report for %s to %s.\n"
				output.Printf(format, start.Format(ReportDateFmt), end.Format(ReportDateFmt))
				output.Println()
			}

			output.Printf("Total Duration: %s\n", duration)
			output.Printf("Entry Count: %d\n", entries)
			output.Println()
		}

		if format != "" {
			// Write formatted output
			return forEachEntry(tsGateway, sheets, func(entry types.Entry) {
				out := reportOutputItem{}
				out.Entry = entry
				out.Status = status

				tmpl := template.Must(template.New("status").Parse(format))
				tmpl.Execute(output.Writer, out)

				// Always end with a new line...
				output.Println()
			})
		}

		// Write table
		table := tablewriter.NewWriter(output.Writer)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		table.SetHeader([]string{
			"Date",
			"Hash",
			"Created",
			"Updated",
			"Note",
			"Duration",
			"Running",
		})

		err = forEachEntry(tsGateway, sheets, func(entry types.Entry) {
			table.Append([]string{
				entry.Timesheet,
				entry.ShortHash(),
				entry.Created.Format(entry.CreatedTimeFormat()),
				entry.Updated.Format(entry.UpdatedTimeFormat()),
				entry.Note,
				entry.Duration.String(),
				fmt.Sprintf("%t", entry.IsRunning),
			})
		})

		if err != nil {
			return err
		}

		table.SetAutoMergeCells(true)
		table.SetRowLine(true)
		table.Render()

		return nil
	}

	return &console.Command{
		Name:        "report",
		Description: "Display a tabular timesheet report.",
		Configure:   configure,
		Execute:     execute,
	}
}

// forEachEntry runs the given function on each entry in each timesheet in the given array of
// timesheets. This uses the database.
func forEachEntry(gw tracking.TimesheetGateway, ss []types.Timesheet, fn func(types.Entry)) error {
	for _, sheet := range ss {
		for _, hash := range sheet.Entries {
			entry, err := gw.FindEntry(hash)
			if err != nil {
				return err
			}

			fn(entry)
		}
	}

	return nil
}

// getTimesheetsByKeys returns all of the timesheets that exist from an array of keys to try.
func getTimesheetsByKeys(gateway tracking.TimesheetGateway, keys []string) ([]types.Timesheet, error) {
	sheets := []types.Timesheet{}

	for _, key := range keys {
		sheet, err := gateway.FindTimesheet(key)
		if err != nil && err != state.ErrStoreNilResult {
			return sheets, err
		}

		if err == state.ErrStoreNilResult {
			continue
		}

		sheets = append(sheets, sheet)
	}

	return sheets, nil
}

// getDateRangeTimesheetKeys produces an array of keys to attempt to find timesheets within for a
// given start and end date range.
func getDateRangeTimesheetKeys(start time.Time, end time.Time) []string {
	keys := []string{}

	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		keys = append(keys, current.Format(types.TimesheetKeyDateFmt))
	}

	return keys
}
