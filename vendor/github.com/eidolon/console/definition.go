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
func (d *Definition) AddArgument(value parameters.Value, spec string, desc string) {
	arg, err := specification.ParseArgumentSpecification(spec)

	if err != nil {
		panic(fmt.Errorf("Error parsing argument specification: '%s'.", err.Error()))
	}

	arg.Value = value
	arg.Description = desc

	if _, ok := d.arguments[arg.Name]; ok {
		panic(fmt.Errorf("Cannot redeclare argument with name '%s'.", arg.Name))
	}

	d.arguments[arg.Name] = arg
	d.argumentKeys = append(d.argumentKeys, arg.Name)
}

// AddOption creates a parameters.Option and adds it to the Definition. Duplicate option names will
// result in an error.
func (d *Definition) AddOption(value parameters.Value, spec string, desc string) {
	opt, err := specification.ParseOptionSpecification(spec)

	if err != nil {
		panic(fmt.Errorf("Error parsing option specification: '%s'.", err.Error()))
	}

	opt.Value = value
	opt.Description = desc

	for _, name := range opt.Names {
		if _, ok := d.options[name]; ok {
			panic(fmt.Errorf("Cannot redeclare option with name '%s'.", name))
		}

		d.options[name] = opt
	}

	d.optionSet = append(d.optionSet, opt)
}
