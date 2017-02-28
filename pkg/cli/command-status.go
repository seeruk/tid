package cli

import (
	"fmt"
	"text/template"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
	"github.com/olekukonko/tablewriter"
)

// statusOutput represents the formattable source of the status command output.
type statusOutput struct {
	Entry   tracking.Entry
	Status  tracking.Status
	Running bool
}

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(gateway tracking.Gateway) console.Command {
	var format string
	var hash string

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewStringValue(&format),
			"-f, --format=FORMAT",
			"Format string, uses table headers e.g. '{{HASH}}'.",
		)

		def.AddArgument(
			parameters.NewStringValue(&hash),
			"[HASH]",
			"A short or long hash for an entry.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.Entry == "" {
			output.Println("status: No timer to check the status of")
			return nil
		}

		if hash == "" {
			hash = status.Entry
		}

		entry, err := gateway.FindEntry(hash)
		if err != nil {
			return err
		}

		dateFormat := "3:04PM (2006-01-02)"
		isRunning := status.IsActive && status.Entry == entry.Hash

		if isRunning {
			// If we're viewing the status of the currently active entry, we should get make sure
			// that it's duration is up-to-date.
			entry.UpdateDuration()
		}

		if format != "" {
			out := statusOutput{}
			out.Entry = entry
			out.Status = status
			out.Running = isRunning

			tmpl := template.Must(template.New("status").Parse(format))
			tmpl.Execute(output.Writer, out)

			// Always end with a new line...
			output.Println()
		} else {
			table := tablewriter.NewWriter(output.Writer)
			table.SetHeader([]string{
				"Date",
				"Hash",
				"Created",
				"Updated",
				"Note",
				"Duration",
				"Running",
			})
			table.Append([]string{
				entry.Timesheet,
				entry.ShortHash(),
				entry.Created.Format(dateFormat),
				entry.Updated.Format(dateFormat),
				entry.Note,
				entry.Duration.String(),
				fmt.Sprintf("%t", isRunning),
			})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.Render()
		}

		return nil
	}

	return console.Command{
		Name:        "status",
		Description: "View the current status.",
		Configure:   configure,
		Execute:     execute,
	}
}
