package parameters_test

import (
	"strings"
	"testing"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestDescribeOptions(t *testing.T) {
	t.Run("should include a title", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{})

		assert.True(t, strings.Contains(result, "OPTIONS:"), "Expected a title.")
	})

	t.Run("should include option names", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names: []string{
					"t",
					"test",
				},
			},
		})

		assert.True(t, strings.Contains(result, "-t"), "Expected option name in result.")
		assert.True(t, strings.Contains(result, "--test"), "Expected option name in result.")
	})

	t.Run("should handle multiple arguments", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names: []string{
					"f",
					"foo",
				},
			},
			{
				Names: []string{
					"b",
					"bar",
				},
			},
		})

		assert.True(t, strings.Contains(result, "-f"), "Expected option name in result.")
		assert.True(t, strings.Contains(result, "--foo"), "Expected option name in result.")
		assert.True(t, strings.Contains(result, "-b"), "Expected option name in result.")
		assert.True(t, strings.Contains(result, "--bar"), "Expected option name in result.")
	})

	t.Run("should sort arguments into alphabetical order", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names: []string{
					"f",
					"foo",
				},
			},
			{
				Names: []string{
					"b",
					"bar",
				},
			},
		})

		fooIdx := strings.Index(result, "--foo")
		barIdx := strings.Index(result, "--bar")

		assert.True(t, fooIdx > barIdx, "Expected FOO to come after BAR.")
	})

	t.Run("should show value names", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names:     []string{"foo"},
				ValueMode: parameters.OptionValueRequired,
				ValueName: "FOO_NAME",
			},
		})

		assert.True(t, strings.Contains(result, "FOO_NAME"), "Expected value name in output.")
	})

	t.Run("should show value names as optional if they are", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names:     []string{"foo"},
				ValueMode: parameters.OptionValueOptional,
				ValueName: "FOO_NAME",
			},
		})

		assert.True(t, strings.Contains(result, "[=FOO_NAME]"), "Expected value name in output.")
	})

	t.Run("should sort short options before long options names", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names: []string{
					"foo",
					"f",
				},
			},
		})

		fooIdx := strings.Index(result, "--foo")
		fIdx := strings.Index(result, "-f")

		assert.True(t, fooIdx > fIdx, "Expected --foo to come after -f.")
	})

	t.Run("should sort short options into alphabetical order", func(t *testing.T) {
		result := parameters.DescribeOptions([]parameters.Option{
			{
				Names: []string{
					"f",
					"b",
				},
			},
		})

		bIdx := strings.Index(result, "-b")
		fIdx := strings.Index(result, "-f")

		assert.True(t, fIdx > bIdx, "Expected -f to come after -b.")
	})
}
