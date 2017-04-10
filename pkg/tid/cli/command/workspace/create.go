package workspace

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// CreateCommand creates a command that creates new workspaces.
func CreateCommand(factory util.Factory) *console.Command {
	var name string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&name),
			Spec:  "NAME",
			Desc:  "A workspace name.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildWorkspaceFacade()

		err := facade.Create(name)
		if err != nil {
			return err
		}

		output.Printf("Created workspace '%s'\n", name)

		return nil
	}

	return &console.Command{
		Name:        "create",
		Alias:       "c",
		Description: "Create a new workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
