package console_test

import (
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestMapInput(t *testing.T) {
	createInput := func(params []string) *console.Input {
		return console.ParseInput(params)
	}

	t.Run("should map arguments to their reference values", func(t *testing.T) {
		var s1 string
		var s2 string

		assert.Equal(t, "", s1)
		assert.Equal(t, "", s2)

		input := createInput([]string{"hello", "world"})

		definition := console.NewDefinition()
		definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
		definition.AddArgument(parameters.NewStringValue(&s2), "S2", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)

		assert.Equal(t, "hello", s1)
		assert.Equal(t, "world", s2)
	})

	t.Run("should error parsing arguments with invalid values", func(t *testing.T) {
		var i1 int

		assert.Equal(t, 0, i1)

		input := createInput([]string{"foo"})

		definition := console.NewDefinition()
		definition.AddArgument(parameters.NewIntValue(&i1), "I1", "")

		err := console.MapInput(definition, input)
		assert.NotOK(t, err)
	})

	t.Run("should error when required arguments are missing from input", func(t *testing.T) {
		var s1 string
		var s2 string

		assert.Equal(t, "", s1)
		assert.Equal(t, "", s2)

		input := createInput([]string{"foo"})

		definition := console.NewDefinition()
		definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
		definition.AddArgument(parameters.NewStringValue(&s1), "S2", "")

		err := console.MapInput(definition, input)
		assert.NotOK(t, err)
	})

	t.Run("should not error when optional arguments are missing from input", func(t *testing.T) {
		var s1 string
		var s2 string

		assert.Equal(t, "", s1)
		assert.Equal(t, "", s2)

		input := createInput([]string{"foo"})

		definition := console.NewDefinition()
		definition.AddArgument(parameters.NewStringValue(&s1), "S1", "")
		definition.AddArgument(parameters.NewStringValue(&s1), "[S2]", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)

		assert.Equal(t, "foo", s1)
		assert.Equal(t, "", s2)
	})

	t.Run("should map short options to their reference values", func(t *testing.T) {
		var s1 string
		var s2 string

		assert.Equal(t, "", s1)
		assert.Equal(t, "", s2)

		input := createInput([]string{"-a=foo", "-b=bar"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewStringValue(&s1), "-a=S1", "")
		definition.AddOption(parameters.NewStringValue(&s2), "-b=S2", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)

		assert.Equal(t, "foo", s1)
		assert.Equal(t, "bar", s2)
	})

	t.Run("should map long options to their reference values", func(t *testing.T) {
		var s1 string
		var s2 string

		assert.Equal(t, "", s1)
		assert.Equal(t, "", s2)

		input := createInput([]string{"--foo=bar", "--baz=qux"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewStringValue(&s1), "--foo=S1", "")
		definition.AddOption(parameters.NewStringValue(&s2), "--baz=S2", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)

		assert.Equal(t, "bar", s1)
		assert.Equal(t, "qux", s2)
	})

	t.Run("should ignore options that don't exist in the definition", func(t *testing.T) {
		var s2 string

		input := createInput([]string{"--foo=bar"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewStringValue(&s2), "--baz=S2", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)
	})

	t.Run("should error parsing an option that requires a value with no value", func(t *testing.T) {
		var s1 string

		input := createInput([]string{"--foo"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewStringValue(&s1), "--foo=s1", "")

		err := console.MapInput(definition, input)
		assert.NotOK(t, err)
	})

	t.Run("should not error parsing an option that doesn't require a value", func(t *testing.T) {
		var s1 string

		input := createInput([]string{"--foo"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewStringValue(&s1), "--foo=s1", "")

		err := console.MapInput(definition, input)
		assert.NotOK(t, err)
	})

	t.Run("should set flag option values where applicable", func(t *testing.T) {
		var b1 bool

		assert.Equal(t, false, b1)

		input := createInput([]string{"--foo"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewBoolValue(&b1), "--foo", "")

		err := console.MapInput(definition, input)
		assert.OK(t, err)

		assert.Equal(t, true, b1)
	})

	t.Run("should error parsing options with invalid values", func(t *testing.T) {
		var i1 int

		assert.Equal(t, 0, i1)

		input := createInput([]string{"--foo=hello"})

		definition := console.NewDefinition()
		definition.AddOption(parameters.NewIntValue(&i1), "--foo=I1", "")

		err := console.MapInput(definition, input)
		assert.NotOK(t, err)
	})
}
