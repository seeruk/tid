package display

import (
	"fmt"
	"io"

	"github.com/SeerUK/tid/pkg/types"
	"github.com/olekukonko/tablewriter"
)

// WriteEntriesTable writes the given entries to a writer as a table.
func WriteEntriesTable(entries []types.Entry, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

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
			entry.Duration.String(),
			fmt.Sprintf("%t", entry.IsRunning),
		})
	}

	table.Render()
}
