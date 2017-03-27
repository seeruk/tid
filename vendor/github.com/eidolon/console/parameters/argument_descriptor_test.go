package parameters_test

import (
	"strings"
	"testing"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestDescribeArguments(t *testing.T) {
	t.Run("should include a title", func(t *testing.T) {
		result := parameters.DescribeArguments([]parameters.Argument{})

		assert.True(t, strings.Contains(result, "ARGUMENTS:"), "Expected a title.")
	})

	t.Run("should include argument names", func(t *testing.T) {
		result := parameters.DescribeArguments([]parameters.Argument{
			{
				Name: "TEST_ARG",
			},
		})

		assert.True(t, strings.Contains(result, "TEST_ARG"), "Expected argument name in result.")
	})

	t.Run("should handle multiple arguments", func(t *testing.T) {
		result := parameters.DescribeArguments([]parameters.Argument{
			{
				Name: "TEST_ARG1",
			},
			{
				Name: "TEST_ARG2",
			},
		})

		assert.True(t, strings.Contains(result, "TEST_ARG1"), "Expected argument name in result.")
		assert.True(t, strings.Contains(result, "TEST_ARG2"), "Expected argument name in result.")
	})

	t.Run("should sort arguments into alphabetical order", func(t *testing.T) {
		result := parameters.DescribeArguments([]parameters.Argument{
			{
				Name: "FOO",
			},
			{
				Name: "BAR",
			},
		})

		fooIdx := strings.Index(result, "FOO")
		barIdx := strings.Index(result, "BAR")

		assert.True(t, fooIdx > barIdx, "Expected FOO to come after BAR.")
	})
}
