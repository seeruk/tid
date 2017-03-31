package specification

import "fmt"

// Parser provides a base implementation for parsing parameter specifications.
type parser struct {
	// The token scanner.
	scanner *Scanner
	// Parser context.
	context struct {
		// The last read token type.
		token Token
		// The last read token value string.
		value string
		// Buffer size.
		size int
	}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *parser) scan() (Token, string) {
	// If we have a token on the buffer, then return it.
	if p.context.size != 0 {
		p.context.size = 0

		return p.context.token, p.context.value
	}

	// Otherwise read the next token from the scanner.
	tok, val := p.scanner.Scan()

	if tok == WS {
		return p.scan()
	}

	// Save it to the buffer in case we unscan later.
	p.context.token, p.context.value = tok, val

	return tok, val
}

// unscan pushes the previously read token back onto the buffer.
func (p *parser) unscan() {
	p.context.size = 1
}

// expectedButActual is a helper for creating parser errors with an expected and actual value.
func (p *parser) expected(expected string, actual string) error {
	return fmt.Errorf("specification: Expected %s, found '%s'", expected, actual)
}

// expectedLen is a helper for creating parser errors for string lengths.
func (p *parser) expectedLen(expected string, expectedLen int, actualLen int) error {
	return fmt.Errorf(
		"specification: Expected %s to be %d character(s), was %d character(s)",
		expected,
		expectedLen,
		actualLen,
	)
}
