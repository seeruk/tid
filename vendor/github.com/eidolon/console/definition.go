package console

import (
	"fmt"

	"github.com/eidolon/console/parameters"
	"github.com/eidolon/console/specification"
)

// Definition represents the collected and configured input parameters of an application.
type Definition struct {
	// Defined arguments for the current application run.
	arguments    map[string]parameters.Argument
	argumentKeys []string

	// Defined options for the current application run.
	options   map[string]parameters.Option
	optionSet []parameters.Option
}

// NewDefinition creates a new Definition with sensible defaults.
func NewDefinition() *Definition {
	definition := Definition{}
	definition.arguments = make(map[string]parameters.Argument)
	definition.options = make(map[string]parameters.Option)

	return &definition
}

// ArgumentDefinition is a struct that represents the entire configuration of a CLI argument.
type ArgumentDefinition struct {
	// The value to reference.
	Value parameters.Value
	// The specification of the argument.
	Spec string
	// The description of the argument.
	Desc string
}

// OptionDefinition is a struct that represents the entire configuration of a CLI option.
type OptionDefinition struct {
	// The value to reference.
	Value parameters.Value
	// The specification of the option.
	Spec string
	// The description of the option.
	Desc string
	// The name of an environment variable to read an option value from.
	EnvVar string
}

// Arguments gets all of the arguments in this Definition.
func (d *Definition) Arguments() []parameters.Argument {
	arguments := []parameters.Argument{}

	for _, key := range d.argumentKeys {
		arguments = append(arguments, d.arguments[key])
	}

	return arguments
}

// Options gets all of the options in this Definition.
func (d *Definition) Options() []parameters.Option {
	return d.optionSet
}

// AddArgument creates a parameters.Argument and adds it to the Definition. Duplicate argument names
// will result in an error.
func (d *Definition) AddArgument(definition ArgumentDefinition) {
	arg, err := specification.ParseArgumentSpecification(definition.Spec)

	if err != nil {
		panic(fmt.Errorf("Error parsing argument specification: '%s'.", err.Error()))
	}

	arg.Description = definition.Desc
	arg.Value = definition.Value

	if _, ok := d.arguments[arg.Name]; ok {
		panic(fmt.Errorf("Cannot redeclare argument with name '%s'.", arg.Name))
	}

	d.arguments[arg.Name] = arg
	d.argumentKeys = append(d.argumentKeys, arg.Name)
}

// AddOption creates a parameters.Option and adds it to the Definition. Duplicate option names will
// result in an error.
func (d *Definition) AddOption(definition OptionDefinition) {
	opt, err := specification.ParseOptionSpecification(definition.Spec)

	if err != nil {
		panic(fmt.Errorf("Error parsing option specification: '%s'.", err.Error()))
	}

	opt.Description = definition.Desc
	opt.EnvVar = definition.EnvVar
	opt.Value = definition.Value

	for _, name := range opt.Names {
		if _, ok := d.options[name]; ok {
			panic(fmt.Errorf("Cannot redeclare option with name '%s'.", name))
		}

		d.options[name] = opt
	}

	d.optionSet = append(d.optionSet, opt)
}
