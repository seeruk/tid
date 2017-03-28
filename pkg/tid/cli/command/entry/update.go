package entry

import (
	"time"

	"github.com/SeerUK/tid/pkg/errhandling"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/SeerUK/tid/pkg/types"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// UpdateCommand creates a command to updated timesheet entries.
func UpdateCommand(factory tracking.Factory) *console.Command {
	var duration time.Duration
	var hash string
	var offset time.Duration
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&hash),
			Spec:  "HASH",
			Desc:  "A short or long hash for an entry.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDurationValue(&duration),
			Spec:  "-d, --duration=DURATION",
			Desc:  "A new duration to set on the entry. Mutually exclusive with offset.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&note),
			Spec:  "-n, --note=NOTE",
			Desc:  "A new note to set on the entry.",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewDurationValue(&offset),
			Spec:  "-o, --offset=OFFSET",
			Desc:  "An offset to modify the duration by (can be negative). Mutually exclusive with duration.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		hasDuration := input.HasOption([]string{"d", "duration"})
		hasNote := input.HasOption([]string{"n", "note"})
		hasOffset := input.HasOption([]string{"o", "offset"})

		if hasDuration && hasOffset {
			output.Println("update: Duration and offset are mutually exclusive")
			return nil
		}

		var entry types.Entry
		var err error

		errs := errhandling.NewErrorStack()
		facade := factory.BuildEntryFacade()

		if hasDuration {
			entry, err = facade.UpdateDuration(hash, duration)
			errs.Add(err)
		}

		if hasOffset {
			entry, err = facade.UpdateDurationByOffset(hash, offset)
			errs.Add(err)
		}

		if hasNote {
			entry, err = facade.UpdateNote(hash, note)
			errs.Add(err)
		}

		if !errs.Empty() {
			return errs.Errors()
		}

		output.Printf("Updated entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "update",
		Description: "Update a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
