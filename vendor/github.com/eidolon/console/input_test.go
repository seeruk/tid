package console_test

import (
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
)

func TestInput(t *testing.T) {
	t.Run("HasOption", func(t *testing.T) {
		t.Run("should return true if a given option exists", func(t *testing.T) {
			input := createTestInput([]string{"example", "e", "foo", "bar"})

			assert.True(t, input.HasOption([]string{"example"}), "Expected option to exist")
			assert.True(t, input.HasOption([]string{"e"}), "Expected option to exist")
			assert.True(t, input.HasOption([]string{"foo"}), "Expected option to exist")
			assert.True(t, input.HasOption([]string{"bar"}), "Expected option to exist")
			assert.True(t, input.HasOption([]string{"e", "example"}), "Expected option to exist")
		})

		t.Run("should return false if a given option doesn't exist", func(t *testing.T) {
			input := createTestInput([]string{})

			assert.False(t, input.HasOption([]string{"example"}), "Expected option not to exist")
			assert.False(t, input.HasOption([]string{"e"}), "Expected option not to exist")
			assert.False(t, input.HasOption([]string{"foo"}), "Expected option not to exist")
			assert.False(t, input.HasOption([]string{"bar"}), "Expected option not to exist")
		})
	})
}

func createTestInput(names []string) console.Input {
	opts := []console.InputOption{}

	for _, name := range names {
		opts = append(opts, console.InputOption{
			Name: name,
		})
	}

	return console.Input{
		Options: opts,
	}
}
