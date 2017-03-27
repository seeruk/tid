package parameters

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eidolon/wordwrap"
)

// DescribeArguments describes an array of Arguments, formatting them in a helpful way.
func DescribeArguments(arguments []Argument) string {
	desc := "ARGUMENTS:\n"

	// Create array and map for specific output ordering
	argDescKeys := []string{}
	argDescMap := make(map[string]string)

	// Generate the list of names and description to allow specific output ordering.
	for _, arg := range arguments {
		argDescKeys = append(argDescKeys, arg.Name)
		argDescMap[arg.Name] = arg.Description
	}

	// Sort option names, so they are output in alphabetical order.
	sort.Sort(argumentNameSort(argDescKeys))

	// Find maximum option names width for spacing.
	var width int
	for _, name := range argDescKeys {
		len := len(name)

		if len+2 > width {
			width = len + 2
		}
	}

	for _, names := range argDescKeys {
		// Get space for the right-side of the command name.
		spacing := width - len(names)

		// Wrap the description onto new lines if necessary.
		wrapper := wordwrap.Wrapper(78-width, true)
		wrapped := wrapper(argDescMap[names])

		// Indent and prefix to product the result.
		prefix := fmt.Sprintf("  %s%s", names, strings.Repeat(" ", spacing))

		desc += wordwrap.Indent(wrapped, prefix, false) + "\n"
	}

	return desc
}

// argumentNameSort allows argument name sorting (trim leading brackets, and alphabetically sort).
type argumentNameSort []string

func (a argumentNameSort) Len() int {
	return len(a)
}

func (a argumentNameSort) Less(i, j int) bool {
	return a[i] < a[j]
}

func (a argumentNameSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
