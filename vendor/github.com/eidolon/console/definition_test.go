package console_test

import (
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestNewDefinition(t *testing.T) {
	definition := console.NewDefinition()
	assert.True(t, definition != nil, "Definition should not be nil")
}

func TestDefinition(t *testing.T) {
	t.Run("Arguments()", func(t *testing.T) {
		t.Run("should return an empty slice if no arguments have been added", func(t *testing.T) {
			definition := console.NewDefinition()
			assert.True(t, len(definition.Arguments()) == 0, "Arguments() length should be 0")
		})

		t.Run("should return an ordered slice of arguments", func(t *testing.T) {
			var s1 string
			var s2 string
			var s3 string

			definition := console.NewDefinition()
			definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
			definition.AddArgument(parameters.NewStringValue(&s2), "S2", "")
			definition.AddArgument(parameters.NewStringValue(&s3), "S3", "")

			arguments := definition.Arguments()

			assert.True(t, len(arguments) == 3, "Arguments() length should be 3")
			assert.Equal(t, "S1", arguments[0].Name)
			assert.Equal(t, "S2", arguments[1].Name)
			assert.Equal(t, "S3", arguments[2].Name)
		})
	})

	t.Run("Options()", func(t *testing.T) {
		t.Run("should return an empty slice if no options have been added", func(t *testing.T) {
			definition := console.NewDefinition()
			assert.True(t, len(definition.Options()) == 0, "Options() length should be 0")
		})

		t.Run("should return a slice of options", func(t *testing.T) {
			var s1 string
			var s2 string
			var s3 string

			definition := console.NewDefinition()
			definition.AddOption(parameters.NewStringValue(&s1), "--s1", "")
			definition.AddOption(parameters.NewStringValue(&s2), "--s2", "")
			definition.AddOption(parameters.NewStringValue(&s3), "--s3", "")

			options := definition.Options()

			assert.Equal(t, 3, len(options))
		})
	})

	t.Run("AddArgument()", func(t *testing.T) {
		t.Run("should error if an invalid option specification is given", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.False(t, r == nil, "We should be recovering from a panic.")
			}()

			var s1 string

			definition := console.NewDefinition()
			definition.AddArgument(parameters.NewStringValue(&s1), "!!! S1", "")
		})

		t.Run("should error if an argument with the same name exists", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.False(t, r == nil, "We should be recovering from a panic.")
			}()

			var s1 string
			var s2 string

			definition := console.NewDefinition()
			definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
			definition.AddArgument(parameters.NewStringValue(&s2), "S1", "")
		})

		t.Run("should add an argument", func(t *testing.T) {
			var s1 string

			definition := console.NewDefinition()
			assert.Equal(t, 0, len(definition.Arguments()))

			definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
			assert.Equal(t, 1, len(definition.Arguments()))
		})
	})

	t.Run("AddOption()", func(t *testing.T) {
		t.Run("should error if an invalid option specification is given", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.False(t, r == nil, "We should be recovering from a panic.")
			}()

			var s1 string

			definition := console.NewDefinition()
			definition.AddOption(parameters.NewStringValue(&s1), "!!! S1", "")
		})

		t.Run("should error if an option with the same name exists", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.False(t, r == nil, "We should be recovering from a panic.")
			}()

			var s1 string
			var s2 string

			definition := console.NewDefinition()
			definition.AddOption(parameters.NewStringValue(&s1), "--s1", "")
			definition.AddOption(parameters.NewStringValue(&s2), "--s1", "")
		})

		t.Run("should add an option", func(t *testing.T) {
			var s1 string

			definition := console.NewDefinition()
			assert.Equal(t, 0, len(definition.Options()))

			definition.AddOption(parameters.NewStringValue(&s1), "--s1", "")
			assert.Equal(t, 1, len(definition.Options()))
		})
	})
}
