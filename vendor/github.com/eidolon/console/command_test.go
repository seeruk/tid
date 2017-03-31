package console_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
)

func TestCommand(t *testing.T) {
	t.Run("AddCommands()", func(t *testing.T) {
		t.Run("should add all of the given sub-commands", func(t *testing.T) {
			inCommands := []*console.Command{
				{
					Name: fmt.Sprintf("test%d", rand.Int()),
				},
				{
					Name: fmt.Sprintf("test%d", rand.Int()),
				},
			}

			command := console.Command{}
			command.AddCommands(inCommands)

			assert.True(t, len(command.Commands()) == 2, "Expected two commands")

			for i, cmd := range command.Commands() {
				assert.Equal(t, inCommands[i], cmd)
			}
		})
	})

	t.Run("AddCommand()", func(t *testing.T) {
		t.Run("should add a given sub-command", func(t *testing.T) {
			inCommand := &console.Command{
				Name: fmt.Sprintf("test%d", rand.Int()),
			}

			command := console.Command{}
			command.AddCommand(inCommand)

			assert.True(t, len(command.Commands()) == 1, "Expected one command")
			assert.Equal(t, inCommand, command.Commands()[0])
		})
	})

	t.Run("Commands()", func(t *testing.T) {
		t.Run("should return all sub-commands", func(t *testing.T) {
			inCommands := []*console.Command{
				{
					Name: fmt.Sprintf("test%d", rand.Int()),
				},
				{
					Name: fmt.Sprintf("test%d", rand.Int()),
				},
			}

			command := console.Command{}
			command.AddCommands(inCommands)

			assert.True(t, len(command.Commands()) == 2, "Expected two commands")

			for i, cmd := range command.Commands() {
				assert.Equal(t, inCommands[i], cmd)
			}
		})
	})
}
