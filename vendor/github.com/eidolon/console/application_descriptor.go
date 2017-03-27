package console

import (
	"fmt"

	"github.com/eidolon/console/parameters"
	"github.com/eidolon/wordwrap"
)

// DescribeApplication describes an Application to provide usage information.
func DescribeApplication(app *Application) string {
	var help string

	if app.Logo != "" {
		help += fmt.Sprintf("%s\n", app.Logo)
	}

	help += fmt.Sprintf("%s version %s\n\n", app.Name, app.Version)
	help += fmt.Sprintf("%s\n", describeApplicationUsage(app))

	options := findApplicationOptions(app)

	if len(options) > 0 {
		help += fmt.Sprintf("\n%s", parameters.DescribeOptions(options))
	}

	if len(app.commands) > 0 {
		help += fmt.Sprintf("\n%s", DescribeCommands(app.commands))
		help += fmt.Sprintf(
			"\n  Run `$ %s COMMAND --help` for more information about a command.\n",
			app.UsageName,
		)
	}

	if len(app.Help) > 0 {
		help += "\nHELP:\n"
		help += wordwrap.Indent(app.Help, "  ", true) + "\n"
	}

	return help
}

// describeApplicationUsage describes the application's usage, mainly based on whether or not the
// application has commands and if it does, whether you must run a command for anything to happen.
func describeApplicationUsage(app *Application) string {
	desc := "USAGE:\n"
	desc += fmt.Sprintf("  %s COMMAND [OPTIONS...] [ARGUMENTS...]", app.UsageName)

	return desc
}

// findApplicationOptions finds all of the defined options on the given application.
func findApplicationOptions(app *Application) []parameters.Option {
	definition := NewDefinition()

	app.preConfigure(definition)

	if app.Configure != nil {
		app.Configure(definition)
	}

	var options []parameters.Option
	for _, opt := range definition.Options() {
		options = append(options, opt)
	}

	return options
}
