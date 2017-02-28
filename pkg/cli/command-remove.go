package cli

import (
	"github.com/SeerUK/tid/pkg/state"
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// RemoveCommand creates a command to remove timesheet entries.
func RemoveCommand(gateway tracking.Gateway) console.Command {
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(
			parameters.NewStringValue(&hash),
			"HASH",
			"A short or long hash for an entry.",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		entry, err := gateway.FindEntry(hash)
		if err != nil && err != state.ErrNilResult {
			return err
		}

		if err == state.ErrNilResult {
			output.Printf("remove: No entry with hash '%s'\n", hash)
			return nil
		}

		err = gateway.RemoveEntry(entry)
		if err != nil {
			return err
		}

		// @todo: Consider adding onSuccess / postExecute to eidolon/console.
		output.Printf("Removed entry '%s' (%s)\n", entry.Note, entry.ShortHash)

		return nil
	}

	return console.Command{
		Name:        "remove",
		Description: "Remove a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
