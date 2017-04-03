package entry

import (
	"time"

	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// CreateCommand creates a command to add timesheet entries.
func CreateCommand(factory util.Factory) *console.Command {
	var duration time.Duration
	var note string
	var started = time.Now()

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewDurationValue(&duration),
			Spec:  "DURATION",
			Desc:  "How long did you spend on what you want to add?",
		})

		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&note),
			Spec:  "NOTE",
			Desc:  "What were you working on?",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDateValue(&started),
			Spec:  "-d, --date=DATE",
			Desc:  "When did you start working? (Default: today)",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildEntryFacade()

		entry, err := facade.Create(started, duration, note)
		if err != nil {
			return err
		}

		output.Printf("Created entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "create",
		Alias:       "c",
		Description: "Create a new timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
