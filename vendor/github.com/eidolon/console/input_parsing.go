package console

import "strings"

// ParseInput takes an array of strings (typically arguments to the application), and parses them
// into the raw Input type.
func ParseInput(params []string) *Input {
	var result Input

	processOptions := true

	for _, param := range params {
		if param == "--" {
			processOptions = false
			continue
		}

		paramLen := len(param)
		isLongOpt := paramLen > 2 && strings.HasPrefix(param, "--")
		isShortOpt := paramLen > 1 && strings.HasPrefix(param, "-")

		if processOptions && isLongOpt {
			result.Options = append(result.Options, parseOption(param, "--"))
		} else if processOptions && isShortOpt {
			result.Options = append(result.Options, parseOption(param, "-"))
		} else {
			result.Arguments = append(result.Arguments, InputArgument{Value: param})
		}
	}

	return &result
}

// parseOption parses an input option with the given prefix (e.g. '-', or '--').
func parseOption(option string, prefix string) InputOption {
	var result InputOption

	trimmed := strings.TrimPrefix(option, prefix)
	split := strings.SplitN(trimmed, "=", 2)

	if len(split) > 1 {
		result = InputOption{Name: split[0], Value: split[1]}
	} else {
		result = InputOption{Name: split[0], Value: ""}
	}

	return result
}
