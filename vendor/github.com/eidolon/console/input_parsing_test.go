package console_test

import (
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
)

func TestParseInput(t *testing.T) {
	t.Run("should return an empty Input if no parameters are given", func(t *testing.T) {
		params := []string{}
		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 0, "Expected length to be 0")
	})

	t.Run("should parse strings not starting with '-' or '--' as arguments", func(t *testing.T) {
		params := []string{
			"hello",
			"world",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 2, "Expected length to be 2")
		assert.True(t, len(input.Options) == 0, "Expected length to be 0")
	})

	t.Run("should retain argument order", func(t *testing.T) {
		params := []string{
			"lorem",
			"ipsum",
			"dolor",
			"sit",
			"amet",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 5, "Expected length to be 5")
		assert.True(t, len(input.Options) == 0, "Expected length to be 0")
		assert.Equal(t, "lorem", input.Arguments[0].Value)
		assert.Equal(t, "ipsum", input.Arguments[1].Value)
		assert.Equal(t, "dolor", input.Arguments[2].Value)
		assert.Equal(t, "sit", input.Arguments[3].Value)
		assert.Equal(t, "amet", input.Arguments[4].Value)
	})

	t.Run("should not parse '--' as a parameter", func(t *testing.T) {
		params := []string{
			"--",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 0, "Expected length to be 0")
	})

	t.Run("should parse short options", func(t *testing.T) {
		params := []string{
			"-a",
			"-b",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 2, "Expected length to be 2")
		assert.Equal(t, "a", input.Options[0].Name)
		assert.Equal(t, "b", input.Options[1].Name)
	})

	t.Run("should parse short options with values", func(t *testing.T) {
		params := []string{
			"-a=foo",
			"-b=bar",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 2, "Expected length to be 2")
		assert.Equal(t, "a", input.Options[0].Name)
		assert.Equal(t, "foo", input.Options[0].Value)
		assert.Equal(t, "b", input.Options[1].Name)
		assert.Equal(t, "bar", input.Options[1].Value)
	})

	t.Run("should parse long options", func(t *testing.T) {
		params := []string{
			"--foo",
			"--bar",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 2, "Expected length to be 2")
		assert.Equal(t, "foo", input.Options[0].Name)
		assert.Equal(t, "bar", input.Options[1].Name)
	})

	t.Run("should parse long options with values", func(t *testing.T) {
		params := []string{
			"--foo=baz",
			"--bar=qux",
		}

		input := console.ParseInput(params)

		assert.True(t, len(input.Arguments) == 0, "Expected length to be 0")
		assert.True(t, len(input.Options) == 2, "Expected length to be 2")
		assert.Equal(t, "foo", input.Options[0].Name)
		assert.Equal(t, "baz", input.Options[0].Value)
		assert.Equal(t, "bar", input.Options[1].Name)
		assert.Equal(t, "qux", input.Options[1].Value)
	})
}
