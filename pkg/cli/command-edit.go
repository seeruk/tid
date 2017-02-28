package cli

import (
	"time"

	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// EditCommand creates a command to edit timesheet entries.
func EditCommand(gateway tracking.Gateway) console.Command {
	var hash string
	var offset time.Duration
	var note string

	// We need an identifiable value for this
	duration := time.Duration(-1)

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewDurationValue(&duration),
			"-d, --duration=DURATION",
			"A new duration to set on the entry. Mutually exclusive with offset.",
		)

		def.AddOption(
			parameters.NewStringValue(&note),
			"-n, --note=NOTE",
			"A new note to set on the entry.",
		)

		def.AddOption(
			parameters.NewDurationValue(&offset),
			"-o, --offset=OFFSET",
			"An offset to modify the duration by (can be negative). Mutually exclusive with duration.",
		)

		def.AddArgument(
			parameters.NewStringValue(&hash),
			"HASH",
			"A short or long hash for an entry.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		// @todo: Maybe a facade?

		entry, err := gateway.FindEntry(hash)
		if err != nil && err != state.ErrNilResult {
			return err
		}

		if err == state.ErrNilResult {
			output.Printf("edit: No entry with hash '%s'\n", hash)
			return nil
		}

		status, err := gateway.FindOrCreateStatus()
		if err != nil {
			return err
		}

		if status.IsActive && status.Entry == entry.Hash {
			entry.UpdateDuration()
		}

		if duration >= 0 && offset != 0 {
			output.Println("edit: Cannot specify duration and offset at the same time")
			return nil
		}

		// Check for "zero-values"
		if duration < 0 {
			duration = entry.Duration
		}

		if note == "" {
			note = entry.Note
		}

		duration = duration + offset

		if duration < 0 {
			output.Println("edit: Duration can not be less than 0")
			return nil
		}

		entry.Duration = duration
		entry.Note = note
		entry.Updated = time.Now()

		err = gateway.PersistEntry(entry)
		if err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Updated entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		// @todo: Pretty table with old vs. new values.

		return nil
	}

	return console.Command{
		Name:        "edit",
		Description: "Edit a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
