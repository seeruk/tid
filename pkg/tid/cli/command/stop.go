package command

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
)

// StopCommand creates a command to stop timers.
func StopCommand(factory util.Factory) *console.Command {
	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildTrackingFacade()

		entry, err := facade.Stop()
		if err != nil {
			return err
		}

		output.Printf("Stopped timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "stop",
		Description: "Stop the current timer.",
		Execute:     execute,
	}
}
