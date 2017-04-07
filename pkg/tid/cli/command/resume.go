package command

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// ResumeCommand creates a command to resume timers.
func ResumeCommand(factory util.Factory) *console.Command {
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&hash),
			Spec:  "[HASH]",
			Desc:  "A short or long hash for an entry.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildTrackingFacade()

		_, err := facade.Stop()
		if err != nil && err != util.ErrNoTimerRunning {
			return err
		}

		entry, err := facade.Resume(hash)
		if err != nil {
			return err
		}

		output.Printf("Resumed timer for '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "resume",
		Alias:       "res",
		Description: "Resume an existing timer.",
		Configure:   configure,
		Execute:     execute,
	}
}
