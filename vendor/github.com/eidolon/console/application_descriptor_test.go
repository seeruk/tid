package console_test

import (
	"strings"
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestDescribeApplication(t *testing.T) {
	t.Run("should show the application logo if there is one", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Logo = "Eidolon Console\n"

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, application.Logo), "Expected application logo.")
	})

	t.Run("should show the application name", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, application.Name), "Expected application name.")
	})

	t.Run("should show the application version", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, application.Version), "Expected application version.")
	})

	t.Run("should show the application usage", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.UsageName = "eidolon_console_binary"

		usage := application.UsageName + " COMMAND [OPTIONS...] [ARGUMENTS...]"

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, "USAGE:"), "Expected application usage title.")
		assert.True(t, strings.Contains(result, usage), "Expected application usage.")
	})

	t.Run("should show the application options if there are any", func(t *testing.T) {
		var s1 string

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Configure = func(definition *console.Definition) {
			definition.AddOption(console.OptionDefinition{
				Value: parameters.NewStringValue(&s1),
				Spec:  "--s1",
				Desc:  "S1 option for testing.",
			})
		}

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, "OPTIONS:"), "Expected application options title.")
		assert.True(t, strings.Contains(result, "-h"), "Expected application options.")
		assert.True(t, strings.Contains(result, "--help"), "Expected application options.")
		assert.True(t, strings.Contains(result, "--s1"), "Expected application options.")
	})

	t.Run("should show the application commands if there are any", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.AddCommands([]*console.Command{
			{
				Name: "foo-cmd",
			},
			{
				Name: "bar-cmd",
			},
		})

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, "COMMANDS:"), "Expected application commands title.")
		assert.True(t, strings.Contains(result, "foo-cmd"), "Expected application commands.")
		assert.True(t, strings.Contains(result, "bar-cmd"), "Expected application commands.")

	})

	t.Run("should not show the commands title if there are no commands", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Logo = "Eidolon Console\n"

		result := console.DescribeApplication(application)

		assert.False(t, strings.Contains(result, "COMMANDS:"), "Expected no commands title.")
	})

	t.Run("should show the application help if there is any", func(t *testing.T) {
		help := "This is some application help right here. Lorem ipsum dolor sit amet."

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Help = help

		result := console.DescribeApplication(application)

		assert.True(t, strings.Contains(result, "HELP:"), "Expected help title.")
		assert.True(t, strings.Contains(result, help), "Expected help.")
	})

	t.Run("should not show the help title if there is no application help", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Logo = "Eidolon Console\n"

		result := console.DescribeApplication(application)

		assert.False(t, strings.Contains(result, "HELP:"), "Expected no help title.")
	})
}
