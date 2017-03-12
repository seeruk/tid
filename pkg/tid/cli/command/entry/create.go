package entry

import (
	"time"

	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// CreateCommand creates a command to add timesheet entries.
func CreateCommand(factory tracking.Factory) *console.Command {
	var started time.Time
	var duration time.Duration
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewDateValue(&started),
			"STARTED",
			"When did you start working?",
		)

		def.AddArgument(
			parameters.NewDurationValue(&duration),
			"DURATION",
			"How long did you spend on what you want to add?",
		)

		def.AddArgument(
			parameters.NewStringValue(&note),
			"NOTE",
			"What were you working on?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildEntryFacade()

		entry, err := facade.Create(started, duration, note)
		if err != nil {
			return err
		}

		output.Printf("Added entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "create",
		Description: "Create a new timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
