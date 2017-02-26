package cli

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(gateway tracking.Gateway) console.Command {
	var short bool
	var hash string

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewBoolValue(&short),
			"-s, --short",
			"Show shortened output?",
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

		if status.Ref() == nil || status.Ref().Entry == "" {
			output.Println("status: No timer to check the status of")
			return nil
		}

		if hash == "" {
			hash = status.Ref().Entry
		}

		entry, err := gateway.FindEntry(hash)
		if err != nil {
			return err
		}

		isRunning := status.IsActive() && status.Ref().Entry == entry.Hash()

		if isRunning {
			// If we're viewing the status of the currently active entry, we should get make sure
			// that it's duration is up-to-date.
			entry.UpdateDuration()
		}

		if short {
			output.Printf("%s on %s\n", entry.Duration(), entry.ShortHash())
		} else {
			dateFormat := "3:04PM (2006-01-02)"

			output.Printf("Hash: %s\n", entry.Hash())
			output.Println()
			output.Printf("Started At: %s\n", entry.Created().Format(dateFormat))
			output.Printf("Last Started At: %s\n", entry.Updated().Format(dateFormat))
			output.Printf("Duration: %s\n", entry.Duration())
			output.Printf("Note: %s\n", entry.Note())
			output.Printf("Running: %t\n", isRunning)
		}

		return nil
	}

	return console.Command{
		Name:        "status",
		Description: "View the current status. What are you tracking?",
		Configure:   configure,
		Execute:     execute,
	}
}
