package cli

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

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
		return nil
	}

	return console.Command{
		Name:        "status",
		Description: "View the current status. What are you tracking?",
		Configure:   configure,
		Execute:     execute,
	}
}
