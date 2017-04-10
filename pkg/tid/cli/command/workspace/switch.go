package workspace

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// SwitchCommand create a command to switch workspaces.
func SwitchCommand(factory util.Factory) *console.Command {
	var name string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&name),
			Spec:  "NAME",
			Desc:  "A workspace name.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		trFacade := factory.BuildTrackingFacade()
		wsFacade := factory.BuildWorkspaceFacade()

		_, err := trFacade.Stop()
		if err != nil && err != util.ErrNoTimerRunning {
			return err
		}

		err = wsFacade.Switch(name)
		if err != nil {
			return err
		}

		output.Printf("Switched to workspace '%s'\n", name)

		return nil
	}

	return &console.Command{
		Name:        "switch",
		Alias:       "s",
		Description: "Switch to another workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
