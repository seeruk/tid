package workspace

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// DeleteCommand creates a command that deletes workspaces.
func DeleteCommand(factory util.Factory) *console.Command {
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

		err := facade.Delete(name)
		if err != nil {
			return err
		}

		output.Printf("Deleted workspace '%s'\n", name)

		return nil
	}

	return &console.Command{
		Name:        "delete",
		Alias:       "d",
		Description: "Delete a workspace.",
		Configure:   configure,
		Execute:     execute,
	}
}
