package console

// CommandContainer is the interface that provides a method to get commands on an object.
type CommandContainer interface {
	// Commands gets commands from an object.
	Commands() []*Command
}

// ConfigureFunc is a function to mutate the input definition to add arguments and options.
type ConfigureFunc func(*Definition)

// ExecuteFunc is a function to perform whatever task this command does.
type ExecuteFunc func(input *Input, output *Output) error

// Command represents a command to run in an application.
type Command struct {
	// The name of the command.
	Name string
	// An optional alias for the name, usually a shortened name.
	Alias string
	// The description of the command.
	Description string
	// Help message for the command.
	Help string
	// Function to configure command-level parameters.
	Configure ConfigureFunc
	// Function to execute when this command is requested.
	Execute ExecuteFunc

	// Array of sub-commands. May contain sub-commands.
	commands []*Command
}

// AddCommands adds sub-commands to the command.
func (c *Command) AddCommands(commands []*Command) *Command {
	c.commands = append(c.commands, commands...)

	return c
}

// AddCommand adds a sub-command to the command.
func (c *Command) AddCommand(command *Command) *Command {
	c.commands = append(c.commands, command)

	return c
}

// Commands gets the sub-commands on a command.
func (c *Command) Commands() []*Command {
	return c.commands
}
