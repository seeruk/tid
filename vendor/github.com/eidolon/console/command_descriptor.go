package console

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eidolon/console/parameters"
	"github.com/eidolon/wordwrap"
)

// DescribeCommand describes a Command on an Application to provide usage information.
func DescribeCommand(app *Application, cmd *Command, path []string) string {
	var help string

	arguments := findCommandArguments(cmd)
	options := findCommandOptions(app, cmd)

	help += fmt.Sprintf("%s\n", describeCommandUsage(app, cmd, arguments, options, path))

	if len(arguments) > 0 {
		help += fmt.Sprintf("\n%s", parameters.DescribeArguments(arguments))
	}

	if len(options) > 0 {
		help += fmt.Sprintf("\n%s", parameters.DescribeOptions(options))
	}

	if len(cmd.commands) > 0 {
		help += fmt.Sprintf("\n%s", DescribeCommands(cmd.commands))
		help += fmt.Sprintf(
			"\n  Run `$ %s %s COMMAND --help` for more information about a command.\n",
			app.UsageName,
			strings.Join(path, " "),
		)
	}

	if len(cmd.Help) > 0 {
		help += "\nHELP:\n"
		help += wordwrap.Indent(cmd.Help, "  ", true) + "\n"
	}

	return help
}

// DescribeCommands describes an array of Commands to provide usage information.
func DescribeCommands(commands []*Command) string {
	desc := "COMMANDS:\n"

	// Create array and map for specific output ordering.
	cmdKeys := []string{}
	cmdMap := make(map[string]*Command)

	var width int
	for _, cmd := range commands {
		cmdKeys = append(cmdKeys, cmd.Name)
		cmdMap[cmd.Name] = cmd

		len := len(cmd.Name)

		if len > (width - 2) {
			width = len + 2
		}
	}

	sort.Strings(cmdKeys)

	for _, name := range cmdKeys {
		cmd := cmdMap[name]

		// Get space for the right-side of the command name.
		spacing := width - len(name)

		// Wrap the description onto new lines if necessary.
		wrapper := wordwrap.Wrapper(78-width, true)
		wrapped := wrapper(cmd.Description)

		// Indent and prefix to produce the result.
		prefix := fmt.Sprintf("  %s%s", cmd.Name, strings.Repeat(" ", spacing))

		desc += wordwrap.Indent(wrapped, prefix, false) + "\n"
	}

	return desc
}

// describeCommandUsage describes a command's usage.
func describeCommandUsage(app *Application, cmd *Command, args []parameters.Argument, opts []parameters.Option, path []string) string {
	desc := "USAGE:\n"
	desc += fmt.Sprintf(
		"  %s %s",
		app.UsageName,
		strings.Join(path, " "),
	)

	if len(opts) > 0 {
		desc += " [OPTIONS...]"
	}

	if len(args) > 0 {
		for _, arg := range args {
			lb := ""
			rb := ""

			if !arg.Required {
				lb = "["
				rb = "]"
			}

			desc += fmt.Sprintf(" %s%s%s", lb, arg.Name, rb)
		}
	}

	if cmd.Description != "" {
		wrapper := wordwrap.Wrapper(78, true)

		desc += "\n\n" + wordwrap.Indent(wrapper(cmd.Description), "  ", true)
	}

	return desc
}

// findCommandArguments finds all of the defined arguments on the given application and command.
func findCommandArguments(cmd *Command) []parameters.Argument {
	definition := NewDefinition()

	if cmd.Configure != nil {
		cmd.Configure(definition)
	}

	var arguments []parameters.Argument
	for _, arg := range definition.Arguments() {
		arguments = append(arguments, arg)
	}

	return arguments
}

// findCommandOptions finds all of the defined options on the given application and command.
func findCommandOptions(app *Application, cmd *Command) []parameters.Option {
	definition := NewDefinition()

	app.preConfigure(definition)

	if app.Configure != nil {
		app.Configure(definition)
	}

	if cmd.Configure != nil {
		cmd.Configure(definition)
	}

	var options []parameters.Option
	for _, opt := range definition.Options() {
		options = append(options, opt)
	}

	return options
}
