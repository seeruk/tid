package cli

import (
	"fmt"
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
	"github.com/olekukonko/tablewriter"
)

const ReportDateFmt = "2006-01-02"

func ReportCommand(gateway tracking.Gateway) console.Command {
	var start time.Time
	var end time.Time

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
	}

	execute := func(input *console.Input, output *console.Output) error {
		hasStart := input.HasOption([]string{"s", "start"})
		hasEnd := input.HasOption([]string{"e", "end"})

		// We need to get the current date, this is a little hacky, but we need it without any time
		now, err := time.Parse(ReportDateFmt, time.Now().Format(ReportDateFmt))
		if err != nil {
			return err
		}

		if !hasStart {
			start = now
		}

		if !hasEnd {
			end = now
		}

		if start.After(end) {
			output.Println("report: The start date must be before the end date")
			return nil
		}

		keys := getDateRangeTimesheetKeys(start, end)
		sheets, err := getTimesheetsBykeys(gateway, keys)
		if err != nil {
			return err
		}

		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		var duration time.Duration
		var entries int

		err = forEachEntry(gateway, sheets, func(entry *tracking.Entry) {
			if status.Ref().Entry == entry.Hash() {
				entry.UpdateDuration()
				entry.Update()

				gateway.PersistEntry(entry)
			}

			duration = duration + entry.Duration()
			entries = entries + 1
		})

		if err != nil {
			return err
		}

		if entries == 0 {
			output.Println("report: No entries within the given time period")
			return nil
		}

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
		output.Printf("Total Entries: %d\n", entries)
		output.Println()

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

		dateFormat := "03:04:05PM (2006-01-02)"

		err = forEachEntry(gateway, sheets, func(entry *tracking.Entry) {
			isRunning := status.IsActive() && status.Ref().Entry == entry.Hash()

			table.Append([]string{
				entry.Timesheet(),
				entry.ShortHash(),
				entry.Created().Format(dateFormat),
				entry.Updated().Format(dateFormat),
				entry.Note(),
				entry.Duration().String(),
				fmt.Sprintf("%t", isRunning),
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

	return console.Command{
		Name:        "report",
		Description: "Display a tabular timesheet report.",
		Configure:   configure,
		Execute:     execute,
	}
}

// forEachEntry runs the given function on each entry in each timesheet in the given array of
// timesheets. This uses the database.
func forEachEntry(gw tracking.Gateway, ss []*tracking.Timesheet, fn func(*tracking.Entry)) error {
	for _, sheet := range ss {
		for _, hash := range sheet.Entries() {
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
func getTimesheetsBykeys(gateway tracking.Gateway, keys []string) ([]*tracking.Timesheet, error) {
	sheets := []*tracking.Timesheet{}

	for _, key := range keys {
		sheet, err := gateway.FindTimesheet(key)
		if err != nil && err != state.ErrNilResult {
			return sheets, err
		}

		if err == state.ErrNilResult {
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
		keys = append(keys, current.Format(tracking.KeyTimesheetDateFmt))
	}

	return keys
}
