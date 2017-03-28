package console_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"fmt"
	"math/rand"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestNewApplication(t *testing.T) {
	application := console.NewApplication("eidolon/console", "1.2.3+testing")
	assert.True(t, application != nil, "Application should not be nil")
}

func TestApplication(t *testing.T) {
	createApplication := func(writer io.Writer) *console.Application {
		application := console.NewApplication("eidolon/console", "1.2.3.+testing")
		application.Writer = writer

		return application
	}

	createTestCommand := func(a *string, b *int) *console.Command {
		return &console.Command{
			Name: "test",
			Configure: func(definition *console.Definition) {
				definition.AddArgument(console.ArgumentDefinition{
					Value: parameters.NewStringValue(a),
					Spec:  "STRINGARG",
				})

				definition.AddOption(console.OptionDefinition{
					Value: parameters.NewIntValue(b),
					Spec:  "--int-opt=VALUE",
				})
			},
			Execute: func(input *console.Input, output *console.Output) error {
				output.Printf("STRINGARG = %s", *a)
				output.Printf("--int-opt = %v", *b)
				return nil
			},
		}
	}

	t.Run("Run()", func(t *testing.T) {
		t.Run("should return exit code 2 if no command was asked for", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)
			code := application.Run([]string{}, []string{})

			assert.Equal(t, 100, code)
		})

		t.Run("should return exit code 2 if no command was found", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)
			code := application.Run([]string{"foo"}, []string{})

			assert.Equal(t, 100, code)
		})

		t.Run("should return exit code 100 if the help flag is set", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)
			code := application.Run([]string{"--help"}, []string{})

			assert.Equal(t, 100, code)
		})

		t.Run("should show application help if the help flag is set", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.Run([]string{"--help"}, []string{})

			output := writer.String()
			containsUsage := strings.Contains(output, "USAGE:")
			containsArguments := strings.Contains(output, "ARGUMENTS:")
			containsOptions := strings.Contains(output, "OPTIONS:")
			containsHelp := containsUsage && containsOptions && !containsArguments

			assert.True(t, containsHelp, "Expected help output.")
		})

		t.Run("should show command help if the help flag is set when running a command", func(t *testing.T) {
			var a string
			var b int

			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.AddCommand(createTestCommand(&a, &b))
			application.Run([]string{"test", "--help"}, []string{})

			output := writer.String()
			containsUsage := strings.Contains(output, "USAGE:")
			containsArguments := strings.Contains(output, "ARGUMENTS:")
			containsOptions := strings.Contains(output, "OPTIONS:")
			containsHelp := containsUsage && containsOptions && containsArguments

			assert.True(t, containsHelp, "Expected help output.")
		})

		t.Run("should return exit code 0 if a command was found, and ran OK", func(t *testing.T) {
			var a string
			var b int

			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.AddCommand(createTestCommand(&a, &b))

			code := application.Run([]string{"test", "aval", "--int-opt=384"}, []string{})

			assert.Equal(t, 0, code)
		})

		t.Run("should return exit code 101 if mapping input fails", func(t *testing.T) {
			var a string
			var b int

			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.AddCommand(createTestCommand(&a, &b))

			code := application.Run([]string{"test", "aval", "--int-opt=hello"}, []string{})

			assert.Equal(t, 101, code)
		})

		t.Run("should return exit code 102 if the command execution fails", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.AddCommand(&console.Command{
				Name: "test",
				Execute: func(input *console.Input, output *console.Output) error {
					return errors.New("Testing errors")
				},
			})

			code := application.Run([]string{"test", "aval", "--int-opt=hello"}, []string{})

			assert.Equal(t, 102, code)
		})

		t.Run("should configure the application definition", func(t *testing.T) {
			var a string
			var b int
			var foo string

			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.Configure = func(definition *console.Definition) {
				definition.AddOption(console.OptionDefinition{
					Value: parameters.NewStringValue(&foo),
					Spec:  "--foo=FOO",
				})
			}

			application.AddCommand(createTestCommand(&a, &b))
			application.Run([]string{"test", "aval", "--foo=bar"}, []string{})

			assert.Equal(t, "bar", foo)
		})

		t.Run("should work with subcommands", func(t *testing.T) {
			message := fmt.Sprintf("Subcommand: %d", rand.Int())

			subcommand := console.Command{
				Name: "subcommand",
				Execute: func(input *console.Input, output *console.Output) error {
					output.Println(message)
					return nil
				},
			}

			command := console.Command{Name: "command"}
			command.AddCommand(&subcommand)

			writer := bytes.Buffer{}
			application := createApplication(&writer)
			application.AddCommand(&command)
			application.Run([]string{"command", "subcommand"}, []string{})

			assert.True(t, strings.Contains(writer.String(), message), "Expected subcommand to run")
		})
	})

	t.Run("AddCommands()", func(t *testing.T) {
		t.Run("should work when adding 1 command", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)

			assert.Equal(t, 0, len(application.Commands()))

			application.AddCommands([]*console.Command{
				{
					Name: "test1",
				},
			})

			assert.Equal(t, 1, len(application.Commands()))
		})

		t.Run("should work when adding no commands", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)

			assert.Equal(t, 0, len(application.Commands()))

			application.AddCommands([]*console.Command{})

			assert.Equal(t, 0, len(application.Commands()))
		})

		t.Run("should work when adding more than 1 command", func(t *testing.T) {
			writer := bytes.Buffer{}
			application := createApplication(&writer)

			assert.Equal(t, 0, len(application.Commands()))

			application.AddCommands([]*console.Command{
				{
					Name: "test1",
				},
				{
					Name: "test2",
				},
				{
					Name: "test3",
				},
			})

			assert.Equal(t, 3, len(application.Commands()))
		})
	})

	t.Run("AddCommand()", func(t *testing.T) {
		writer := bytes.Buffer{}
		application := createApplication(&writer)

		assert.Equal(t, 0, len(application.Commands()))

		application.AddCommand(&console.Command{
			Name: "test1",
		})

		assert.Equal(t, 1, len(application.Commands()))
	})
}
