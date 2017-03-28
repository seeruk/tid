package console

import (
	"fmt"
	"strings"

	"github.com/eidolon/console/parameters"
)

// MapInput maps the values of input to their corresponding reference values.
func MapInput(definition *Definition, input *Input, env []string) error {
	if err := mapArguments(definition.Arguments(), input); err != nil {
		return err
	}

	if err := mapOptions(definition.Options(), input); err != nil {
		return err
	}

	if err := mapEnv(definition.Options(), env); err != nil {
		return err
	}

	return nil
}

// mapArguments maps the values of input arguments to their corresponding references.
func mapArguments(args []parameters.Argument, input *Input) error {
	var unmappedArguments []parameters.Argument

	for i, arg := range args {
		if len(input.Arguments) == i {
			unmappedArguments = append(unmappedArguments, args[i:]...)
			break
		}

		value := input.Arguments[i].Value

		if err := arg.Value.Set(value); err != nil {
			return fmt.Errorf("Invalid value '%s' for argument '%s'. Error: %s.", value, arg.Name, err)
		}
	}

	for _, uarg := range unmappedArguments {
		if uarg.Required {
			return fmt.Errorf("Argument '%s' is required.", uarg.Name)
		}
	}

	return nil
}

// mapOptions maps the values of input options to their corresponding references.
func mapOptions(opts []parameters.Option, input *Input) error {
	for _, opt := range opts {
		inputOpt := findOptionInInput(opt, input)

		if inputOpt == nil {
			// Option not found in input
			continue
		}

		err := setOptionValue(opt, inputOpt.Name, inputOpt.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// mapEnv maps the values of environment variables into their corresponding option references.
func mapEnv(opts []parameters.Option, env []string) error {
	envMap := make(map[string]string)

	// Split array of option key and values into map.
	for _, ev := range env {
		pair := strings.Split(ev, "=")

		envMap[pair[0]] = pair[1]
	}

	for _, opt := range opts {
		name := ""
		value := ""

		for ek, ev := range envMap {
			if ek == opt.EnvVar {
				name = ek
				value = ev
			}
		}

		if name == "" {
			continue
		}

		err := setOptionValue(opt, name, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// setOptionValue sets the value of an option, and handles potential error cases.
func setOptionValue(opt parameters.Option, name string, value string) error {
	if opt.ValueMode == parameters.OptionValueRequired && value == "" {
		return fmt.Errorf("Option '%s' requires a value.", name)
	}

	isEmptyOptional := opt.ValueMode == parameters.OptionValueOptional && value == ""

	// If we have a flag option, and we received no value, then we should use the preset flag
	// value for if the flag is present.
	if ov, ok := opt.Value.(parameters.FlagValue); value == "" && ok {
		ov.Set(ov.FlagValue())
	} else if !isEmptyOptional {
		err := opt.Value.Set(value)
		if err != nil {
			return fmt.Errorf("Invalid value '%s' for option '%s'. Error: %s.", value, name, err)
		}
	}

	return nil
}

// findOptionInInput finds a given option in the given parsed raw input.
func findOptionInInput(opt parameters.Option, input *Input) *InputOption {
	inputOptions := make(map[string]InputOption)

	for _, inputOption := range input.Options {
		inputOptions[inputOption.Name] = inputOption
	}

	for _, name := range opt.Names {
		if value, ok := inputOptions[name]; ok {
			return &value
		}
	}

	return nil
}
