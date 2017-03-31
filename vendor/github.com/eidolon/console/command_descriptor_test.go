package console_test

import (
	"strings"
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestDescribeCommand(t *testing.T) {
	t.Run("should return command usage information", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name: "test",
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "USAGE:"), "Expected usage information.")
	})

	t.Run("should include the command name", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name: "test-command-name",
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "test-command-name"), "Expected command name.")
	})

	t.Run("should show the command description", func(t *testing.T) {
		description := "This is the test-command-name description."

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name:        "test-command-name",
			Description: description,
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, description), "Expected command description.")
	})

	t.Run("should return command help if there is some", func(t *testing.T) {
		help := "This is some help for the test-command-name command."

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name: "test-command-name",
			Help: help,
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "HELP:"), "Expected command help header.")
		assert.True(t, strings.Contains(result, help), "Expected command help.")
	})

	t.Run("should show arguments if there are any", func(t *testing.T) {
		var s1 string
		var s2 string

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name: "test-command-name",
			Configure: func(definition *console.Definition) {
				definition.AddArgument(console.ArgumentDefinition{
					Value: parameters.NewStringValue(&s1),
					Spec:  "STRING_ARG_S1",
				})

				definition.AddArgument(console.ArgumentDefinition{
					Value: parameters.NewStringValue(&s2),
					Spec:  "STRING_ARG_S2",
				})
			},
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "STRING_ARG_S1"), "Expected argument name.")
		assert.True(t, strings.Contains(result, "STRING_ARG_S2"), "Expected argument name.")
	})

	t.Run("should show optional arguments wrapped in brackets", func(t *testing.T) {
		var s1 string
		var s2 string

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		command := console.Command{
			Name: "test-command-name",
			Configure: func(definition *console.Definition) {
				definition.AddArgument(console.ArgumentDefinition{
					Value: parameters.NewStringValue(&s1),
					Spec:  "[STRING_ARG_S1]",
				})

				definition.AddArgument(console.ArgumentDefinition{
					Value: parameters.NewStringValue(&s2),
					Spec:  "[STRING_ARG_S2]",
				})
			},
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "[STRING_ARG_S1]"), "Expected argument name.")
		assert.True(t, strings.Contains(result, "[STRING_ARG_S2]"), "Expected argument name.")
	})

	t.Run("should show that there are options if there are any", func(t *testing.T) {
		var s1 string
		var s2 string

		application := console.NewApplication("eidolon/console", "1.2.3+testing")
		application.Configure = func(definition *console.Definition) {
			definition.AddOption(console.OptionDefinition{
				Value: parameters.NewStringValue(&s1),
				Spec:  "--s1=VALUE",
			})
		}

		command := console.Command{
			Name: "test-command-name",
			Configure: func(definition *console.Definition) {
				definition.AddOption(console.OptionDefinition{
					Value: parameters.NewStringValue(&s2),
					Spec:  "--s2=VALUE",
				})
			},
		}

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "[OPTIONS...]"), "Expected options.")
	})

	t.Run("should show that there are sub-commands if there are any", func(t *testing.T) {
		application := console.NewApplication("eidolon/console", "1.2.3+testing")

		command := console.Command{
			Name: "test-command-name",
		}

		subCommand := console.Command{
			Name:        "test-subCommand-name",
			Description: "Test sub-command description.",
		}

		command.AddCommand(&subCommand)

		result := console.DescribeCommand(application, &command, []string{command.Name})

		assert.True(t, strings.Contains(result, "COMMANDS:"), "Expected commands")
		assert.True(t, strings.Contains(result, subCommand.Name), "Expected sub-command name")
		assert.True(t, strings.Contains(result, subCommand.Description), "Expected sub-command desc")
	})
}

func TestDescribeCommands(t *testing.T) {
	t.Run("should include a title", func(t *testing.T) {
		result := console.DescribeCommands([]*console.Command{})

		assert.True(t, strings.Contains(result, "COMMANDS:"), "Expected title.")
	})

	t.Run("should show all command names", func(t *testing.T) {
		result := console.DescribeCommands([]*console.Command{
			{
				Name: "foo-cmd",
			},
			{
				Name: "bar-cmd",
			},
		})

		assert.True(t, strings.Contains(result, "foo-cmd"), "Expected command name.")
		assert.True(t, strings.Contains(result, "bar-cmd"), "Expected command name.")
	})

	t.Run("should show all command descriptions", func(t *testing.T) {
		fooCmdDesc := "The foo-cmd description is this."
		barCmdDesc := "An alternative command description for bar-cmd."

		result := console.DescribeCommands([]*console.Command{
			{
				Name:        "foo-cmd",
				Description: fooCmdDesc,
			},
			{
				Name:        "bar-cmd",
				Description: barCmdDesc,
			},
		})

		assert.True(t, strings.Contains(result, fooCmdDesc), "Expected command description.")
		assert.True(t, strings.Contains(result, barCmdDesc), "Expected command description.")
	})

	t.Run("should show the commands in alphabetical order", func(t *testing.T) {
		result := console.DescribeCommands([]*console.Command{
			{
				Name: "foo-cmd",
			},
			{
				Name: "bar-cmd",
			},
		})

		fooIdx := strings.Index(result, "foo-cmd")
		barIdx := strings.Index(result, "bar-cmd")

		assert.True(t, fooIdx > barIdx, "Expected foo-cmd to come after bar-cmd.")
	})
}
