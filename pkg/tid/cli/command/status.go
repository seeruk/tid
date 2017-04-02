package command

import (
	"text/template"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tid/cli/display"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(factory util.Factory) *console.Command {
	var format string
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&hash),
			Spec:  "[HASH]",
			Desc:  "A short or long hash for an entry.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&format),
			Spec:  "-f, --format=FORMAT",
			Desc:  "Output formatting string. Uses Go templates.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		sysGateway := factory.BuildSysGateway()
		trGateway := factory.BuildTrackingGateway()

		hasFormat := input.HasOption([]string{"f", "format"})

		status, err := sysGateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if hash == "" {
			hash = status.Entry
		}

		if hash == "" {
			output.Println("status: No timer to check the status of")
			return nil
		}

		entry, err := trGateway.FindEntry(hash)
		if err != nil && err != state.ErrStoreNilResult {
			return err
		}

		if err == state.ErrStoreNilResult {
			output.Printf("status: No entry with hash '%s'\n", hash)
			return nil
		}

		if hasFormat {
			tmpl := template.Must(template.New("entry-list").Parse(format))
			tmpl.Execute(output.Writer, entry)

			output.Println()

			return nil
		}

		display.WriteEntriesTable([]types.Entry{entry}, output.Writer)

		return nil
	}

	return &console.Command{
		Name:        "status",
		Alias:       "st",
		Description: "View the current status.",
		Configure:   configure,
		Execute:     execute,
	}
}
