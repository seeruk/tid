package console

// Input represents the raw application input, in a slightly more organised way, and provides
// helpers for retrieving that information. This allows commands to get application-wide input,
// instead of only the command-specific input.
type Input struct {
	Arguments []InputArgument
	Options   []InputOption
}

// InputArgument represents the raw data parsed as arguments, really this is just the value.
type InputArgument struct {
	Value string
}

// InputOption represents the raw data parsed as options. This includes it's name and it's value.
// The value can be "" if no value is given (i.e. if the option is a flag).
type InputOption struct {
	Name  string
	Value string
}

// GetOptionValue gets the an option with one of the given names' value./
func (i *Input) GetOptionValue(names []string) string {
	for _, name := range names {
		for _, option := range i.Options {
			if option.Name == name {
				return option.Value
			}
		}
	}

	return ""
}

// HasOption checks to see if the given option exists by one of it's names.
func (i *Input) HasOption(names []string) bool {
	for _, name := range names {
		for _, option := range i.Options {
			if option.Name == name {
				return true
			}
		}
	}

	return false
}

// @todo: Input (low priority):
// @todo: - Add method for retrieving argument by index.
