package command

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
)

// StopCommand creates a command to stop timers.
func StopCommand(factory tracking.Factory) *console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildFacade()

		entry, err := facade.Stop()
		if err != nil {
			return err
		}

		output.Printf("Stopped timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "stop",
		Description: "Stop an existing timer.",
		Execute:     execute,
	}
}
