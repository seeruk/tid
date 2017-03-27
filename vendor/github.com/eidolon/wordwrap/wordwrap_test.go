package wordwrap_test

import (
	"testing"

	"github.com/eidolon/wordwrap"
)

func TestIndent(t *testing.T) {
	tests := []struct {
		input     string
		prefix    string
		prefixAll bool
		expected  string
	}{
		// When testing with no input
		{
			"",
			"",
			false,
			"",
		},
		// When not prefixing all lines
		// Should apply the prefix to the first line
		{
			"Test text\nTest text\nTest text",
			"First line",
			false,
			"First lineTest text\n          Test text\n          Test text",
		},
		// Should allow prefixes to simply be spaces
		{
			"Test text\nTest text",
			"  ",
			false,
			"  Test text\n  Test text",
		},
		// When prefixing all lines
		// Should apply the prefix to all lines
		{
			"Test text\nTest text\nTest text",
			"First line",
			true,
			"First lineTest text\nFirst lineTest text\nFirst lineTest text",
		},
	}

	for i, test := range tests {
		actual := wordwrap.Indent(test.input, test.prefix, test.prefixAll)

		if actual != test.expected {
			t.Fatalf(
				"Result for case %d did not match expected result.\nExpected:\n%s\nActual:\n%s\n",
				i+1,
				test.expected,
				actual,
			)
		}
	}
}

func TestWrapper(t *testing.T) {
	tests := []struct {
		limit      int
		breakWords bool
		input      string
		expected   string
	}{
		// When testing with no input
		{
			4,
			false,
			"",
			"",
		},
		// When not breaking words
		// Should wrap text so words fit on lines with the given limit
		{
			10,
			false,
			"Test text Test text Test text Test text",
			"Test text\nTest text\nTest text\nTest text",
		},
		// Should remove additional whitespace
		{
			10,
			false,
			"Test  text  Test  text  Test  text  Test  text",
			"Test text\nTest text\nTest text\nTest text",
		},
		// Should trim text
		{
			10,
			false,
			"    Test  text  Test  text  Test  text  Test  text    ",
			"Test text\nTest text\nTest text\nTest text",
		},
		// Should not break words if breakWords is false
		{
			4,
			false,
			"Testtext",
			"Testtext",
		},
		// When breaking words
		// Should break words, and insert a hyphen
		{
			4,
			true,
			"Testtext",
			"Test\ntext",
		},
		// Should break words into tiny pieces
		{
			1,
			true,
			"Testtext",
			"T\ne\ns\nt\nt\ne\nx\nt",
		},
		// Should not break words if no words are too long
		{
			4,
			true,
			"Test text",
			"Test\ntext",
		},
		// When given slightly more realistic data
		// Should be broken properly
		{
			40,
			false,
			"This is a bunch of text that should be split at somewhere near 40 characters.",
			"This is a bunch of text that should be\nsplit at somewhere near 40 characters.",
		},
		// When given text with a line break in it
		// It should rebuild the string with spaces in, remember, we're intentionally being simple
		{
			20,
			true,
			"Test\n\ntext",
			"Test text",
		},
	}

	for i, test := range tests {
		wrapper := wordwrap.Wrapper(test.limit, test.breakWords)
		actual := wrapper(test.input)

		if actual != test.expected {
			t.Fatalf(
				"Result for case %d did not match expected result.\nExpected:\n%s\nActual:\n%s\n",
				i+1,
				test.expected,
				actual,
			)
		}
	}
}

func TestWrapperPanicsWithInvalidLimit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Wrapper did not panic on invalid input.")
		}
	}()

	// 0 is an invalid limit
	wordwrap.Wrapper(0, false)
}
