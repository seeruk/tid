package command

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// StartCommand creates a command to start timers.
func StartCommand(factory util.Factory) *console.Command {
	var note string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&note),
			Spec:  "NOTE",
			Desc:  "What are you working on?",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildTrackingFacade()

		entry, err := facade.Start(note)
		if err != nil {
			return err
		}

		output.Printf("Started timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "start",
		Description: "Start a new timer.",
		Configure:   configure,
		Execute:     execute,
	}
}
