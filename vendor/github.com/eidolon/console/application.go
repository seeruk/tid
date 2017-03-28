package console

import (
	"io"
	"os"
	"path/filepath"

	"github.com/eidolon/console/parameters"
)

// Application represents the heart of the console application. It is what orchestrates running
// commands, initiates input parsing, mapping, and validation; and will handle failure for each of
// those tasks.
type Application struct {
	// The name of the application.
	Name string
	// The name of the application printed in usage information, defaults to the binary's filename.
	UsageName string
	// The version of the application.
	Version string
	// Application logo, shown in help output
	Logo string
	// Help message for the application.
	Help string
	// Function to configure application-level parameters, realistically should just be options.
	Configure ConfigureFunc
	// Writer to write output to.
	Writer io.Writer

	// Array of commands that can be run. May contain sub-commands.
	commands []*Command
	// Application input.
	input *Input
	// The path taken to reach the current command (used for help text).
	path []string
}

// NewApplication creates a new Application with some sane defaults.
func NewApplication(name string, version string) *Application {
	return &Application{
		Name:      name,
		UsageName: filepath.Base(os.Args[0]),
		Version:   version,
		Writer:    os.Stdout,
	}
}

// Run runs the configured application, with the given input.
func (a *Application) Run(params []string, env []string) int {
	// Create input and output.
	input := ParseInput(params)
	output := NewOutput(a.Writer)
	definition := NewDefinition()

	// Assign input to application.
	a.input = input
	a.preConfigure(definition)

	if a.Configure != nil {
		// Similar to a.preConfigure, this is useful for customising application-level help text.
		// This however will work just like a command's Configure function.
		a.Configure(definition)

		// Throw away any defined arguments, as they cannot be used, and will disrupt help text.
		definition.arguments = make(map[string]parameters.Argument)
	}

	command, path := a.findCommandInInput()
	if command != nil && command.Configure != nil {
		command.Configure(definition)
	}

	if a.hasHelpOption() || (command == nil || command.Execute == nil) {
		a.showHelp(output, command, path)
		return 100
	}

	err := MapInput(definition, input, env)
	if err != nil {
		output.Println(err)
		output.Printf("Try '%s --help' for more information.\n", a.UsageName)
		return 101
	}

	err = command.Execute(input, output)
	if err != nil {
		output.Println(err)
		output.Printf("Try '%s %s --help' for more information.\n", a.UsageName, command.Name)
		return 102
	}

	return 0
}

// AddCommands adds commands to the application.
func (a *Application) AddCommands(commands []*Command) {
	a.commands = append(a.commands, commands...)
}

// AddCommand adds a command to the application.
func (a *Application) AddCommand(command *Command) {
	a.commands = append(a.commands, command)
}

// Commands gets the subcommands on an application.
func (a *Application) Commands() []*Command {
	return a.commands
}

// findCommandInInput attempts to find the command to run based on the raw input.
func (a *Application) findCommandInInput() (*Command, []string) {
	var loop func(depth int, container CommandContainer) *Command
	var path []string

	loop = func(depth int, container CommandContainer) *Command {
		if len(a.input.Arguments) == 0 {
			return nil
		}

		var command *Command
		for _, cmd := range container.Commands() {
			if cmd.Name == a.input.Arguments[0].Value {
				command = cmd
				// Add to breadcrumb trail...
				path = append(path, cmd.Name)
				break
			}
		}

		if command != nil {
			a.input.Arguments = a.input.Arguments[1:]

			subcommand := loop(depth+1, command)

			if subcommand != nil {
				command = subcommand
			}
		}

		return command
	}

	return loop(0, a), path
}

// hasHelpOption checks to see if a help flag is set, ignoring values. Uses raw input, not mapped
// input.
func (a *Application) hasHelpOption() bool {
	for _, opt := range a.input.Options {
		if opt.Name == "help" || opt.Name == "h" {
			return true
		}
	}

	return false
}

// preConfigure configures pre-defined parameters. This is solely defined for help output.
func (a *Application) preConfigure(definition *Definition) {
	var help bool

	definition.AddOption(OptionDefinition{
		Value: parameters.NewBoolValue(&help),
		Spec:  "-h, --help",
		Desc:  "Display contextual help?",
	})
}

// showHelp shows contextual help.
func (a *Application) showHelp(output *Output, command *Command, path []string) {
	if command != nil {
		output.Println(DescribeCommand(a, command, path))
	} else {
		output.Println(DescribeApplication(a))
	}
}
