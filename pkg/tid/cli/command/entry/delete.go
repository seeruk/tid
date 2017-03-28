package entry

import (
	"github.com/SeerUK/tid/pkg/tracking"
	"github.com/eidolon/console"
	"github.com/eidolon/console/parameters"
)

// DeleteCommand creates a command that is used to delete timesheet entries.
func DeleteCommand(factory tracking.Factory) *console.Command {
	var hash string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&hash),
			Spec:  "HASH",
			Desc:  "A short or long hash for an entry.",
		})
	}

	execute := func(input *console.Input, output *console.Output) error {
		facade := factory.BuildEntryFacade()

		entry, err := facade.Delete(hash)
		if err != nil {
			return err
		}

		output.Printf("Deleted entry '%s' (%s)\n", entry.Note, entry.ShortHash())

		return nil
	}

	return &console.Command{
		Name:        "delete",
		Description: "Delete a timesheet entry.",
		Configure:   configure,
		Execute:     execute,
	}
}
