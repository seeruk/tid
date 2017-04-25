package display

import (
	"fmt"
	"io"
	"time"

	"github.com/SeerUK/tid/pkg/types"
	"github.com/olekukonko/tablewriter"
	"strconv"
)

// WriteEntriesTable writes the given entries to a writer as a table.
func WriteEntriesTable(entries []types.Entry, writer io.Writer, config types.TomlConfig) {
	table := createTable(writer)
	table.SetHeader([]string{
		"Date",
		"Hash",
		"Created",
		"Updated",
		"Note",
		"Duration",
		"Running",
	})

	for _, entry := range entries {
		table.Append([]string{
			entry.Timesheet,
			entry.ShortHash(),
			entry.Created.Format(entry.CreatedTimeFormat()),
			entry.Updated.Format(entry.UpdatedTimeFormat()),
			entry.Note,
			getTimeInRightFormat(entry.Duration, config.Display.TimeFormat),
			fmt.Sprintf("%t", entry.IsRunning),
		})
	}

	table.Render()
}

// WriteTimesheetsTable writes the given timesheets to a writer as a table.
func WriteTimesheetsTable(sheets []types.Timesheet, writer io.Writer, config types.TomlConfig) {
	table := createTable(writer)
	table.SetHeader([]string{
		"Date",
		"Entries",
		"Duration",
	})

	var totalDuration time.Duration
	var totalEntries int

	for _, sheet := range sheets {
		var duration time.Duration

		for _, e := range sheet.Entries {
			duration = duration + e.Duration
		}

		totalDuration = totalDuration + duration
		totalEntries = totalEntries + len(sheet.Entries)

		table.Append([]string{
			sheet.Key,
			fmt.Sprintf("%d", len(sheet.Entries)),
			getTimeInRightFormat(duration, config.Display.TimeFormat),
		})
	}

	// Footer, without affecting value formats
	table.Append([]string{
		"TOTAL",
		fmt.Sprintf("%d", totalEntries),
		totalDuration.String(),
	})

	table.Render()
}

// getTimeInRightFormat prints the time using the format specified in config file.
func getTimeInRightFormat(duration time.Duration, timeFormat string) string {
	if (timeFormat == "decimal") {
		return 	strconv.FormatFloat(duration.Hours(), 'f', 2, 64)
	}

	return duration.String()
}

// createTable creates the base table instance with some default options set.
func createTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	return table
}
