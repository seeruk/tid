package parameters

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eidolon/wordwrap"
)

// DescribeOptions describes an array of Options, formatting them in a helpful way.
func DescribeOptions(options []Option) string {
	desc := "OPTIONS:\n"

	// Create array and map for specific output ordering
	optDescKeys := []string{}
	optDescMap := make(map[string]string)

	// Generate the list of names, so that output can be formatted correctly, and in the correct
	// order (i.e. alphabetical), and with sorted names for each option individually.
	for _, opt := range options {
		var names []string
		for _, name := range opt.Names {
			if len(name) > 1 {
				name = "--" + name
			} else {
				name = "-" + name
			}

			names = append(names, name)
		}

		// Sort the names so that short names appear first in the output.
		sort.Sort(stringLengthSort(names))

		// Join the names into on comma-separated string.
		key := strings.Join(names, ", ")

		// Describe option value
		if opt.ValueMode == OptionValueOptional {
			key += "["
		}

		if opt.ValueMode == OptionValueOptional || opt.ValueMode == OptionValueRequired {
			key += "=" + opt.ValueName
		}

		if opt.ValueMode == OptionValueOptional {
			key += "]"
		}

		optDescKeys = append(optDescKeys, key)
		optDescMap[key] = opt.Description
	}

	// Sort option names, so they are output in alphabetical order.
	sort.Sort(optionNameSort(optDescKeys))

	// Find maximum option names width for spacing.
	var width int
	for _, names := range optDescKeys {
		len := len(names)

		if len+2 > width {
			width = len + 2
		}
	}

	for _, names := range optDescKeys {
		// Get space for the right-side of the command name.
		spacing := width - len(names)

		// Wrap the description onto new lines if necessary.
		wrapper := wordwrap.Wrapper(78-width, true)
		wrapped := wrapper(optDescMap[names])

		// Indent and prefix to product the result.
		prefix := fmt.Sprintf("  %s%s", names, strings.Repeat(" ", spacing))

		desc += wordwrap.Indent(wrapped, prefix, false) + "\n"
	}

	return desc
}

// optionNameSort allows option name sorting (trim leading hyphens, and alphabetically sort).
type optionNameSort []string

func (a optionNameSort) Len() int {
	return len(a)
}

func (a optionNameSort) Less(i, j int) bool {
	l := strings.TrimLeft(a[i], "-")
	r := strings.TrimLeft(a[j], "-")

	return l < r
}

func (a optionNameSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// stringLengthSort allows string length sorting, first tries to sort by length, falls back to
// alphabetical sorting.
type stringLengthSort []string

func (a stringLengthSort) Len() int {
	return len(a)
}

func (a stringLengthSort) Less(i, j int) bool {
	if len(a[i]) == len(a[j]) {
		return a[i] < a[j]
	}

	return len(a[i]) < len(a[j])
}

func (a stringLengthSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
