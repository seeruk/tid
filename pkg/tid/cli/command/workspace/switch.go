package workspace

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// SwitchCommand create a command to switch workspaces.
func SwitchCommand(factory util.Factory) *console.Command {
	var workspace string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&workspace),
			Spec:  "WORKSPACE",
			Desc:  "A workspace name.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildWorkspaceFacade()

		err := facade.Switch(workspace)
		if err != nil {
			return err
		}

		output.Printf("Switched to workspace '%s'\n", workspace)

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
