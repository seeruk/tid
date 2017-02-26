package cli

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StatusCommand creates a command to view the status of the current timer.
func StatusCommand(gateway tracking.Gateway) console.Command {
	var short bool

	configure := func(def *console.Definition) {
		def.AddOption(
			parameters.NewBoolValue(&short),
			"-s, --short",
			"Show shortened output?",
		)
	}

	execute := func(input *console.Input, output *console.Output) error {
		status, err := gateway.FindStatus()
		if err != nil {
			return err
		}

		if !status.IsActive() {
			output.Println("status: There is no active timer running")
			return nil
		}

		entry, err := gateway.FindEntry(status.Ref().Entry)
		if err != nil {
			return err
		}

		// Update the duration, we're not persisting it in this command though.
		entry.UpdateDuration()

		if short {
			output.Printf("%s on %s\n", entry.Duration(), entry.ShortKey())
		} else {
			dateFormat := "3:04PM (2006-01-02)"

			output.Printf("Started At: %s\n", entry.Created().Format(dateFormat))
			output.Printf("Last Started At: %s\n", entry.Updated().Format(dateFormat))
			output.Printf("Duration: %s\n", entry.Duration())
			output.Printf("Note: %s\n", entry.Note())
			output.Println()
			output.Printf("Hash: %s\n", entry.Key())
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
