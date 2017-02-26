package cli

import (
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
		return nil
	}

	return console.Command{
		Name:        "remove",
		Description: "Remove a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
