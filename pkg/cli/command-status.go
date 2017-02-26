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
			output.Printf("%s on %s\n", entry.Duration(), entry.Note())
		} else {
			// @todo: When pausing is in should this show the different start times that there have
			// been (or at least a friendly way of showing that, like a timeline type thing?) Or at
			// the very least, the number of times that it has been paused?
			output.Printf("Started At: %s\n", entry.Created().Format("3:04PM (2006-01-02)"))
			output.Printf("Duration: %s\n", entry.Duration())
			output.Printf("Note: %s\n", entry.Note())
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
