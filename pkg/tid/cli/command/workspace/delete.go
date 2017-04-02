package workspace

import (
	"github.com/SeerUK/tid/pkg/util"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// DeleteCommand creates a command that deletes workspaces.
func DeleteCommand(factory util.Factory) *console.Command {
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

		err := facade.Delete(workspace)
		if err != nil {
			return err
		}

		output.Printf("Deleted workspace '%s'\n", workspace)

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
