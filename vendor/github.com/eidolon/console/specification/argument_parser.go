package specification

import (
	"strings"

	"github.com/eidolon/console/parameters"
)

// ParseArgumentSpecification parses an argument spec string and produces an Argument.
func ParseArgumentSpecification(spec string) (parameters.Argument, error) {
	scanner := NewScanner(strings.NewReader(spec))
	parser := newArgumentSpecifcationParser(scanner)

	return parser.parse()
}

// argumentSpecifcationParser parses argument specification strings.
type argumentSpecifcationParser struct {
	parser
}

// newArgumentSpecifcationParser creates a new argument specification parser.
func newArgumentSpecifcationParser(scanner *Scanner) *argumentSpecifcationParser {
	parser := parser{
		scanner: scanner,
	}

	return &argumentSpecifcationParser{
		parser: parser,
	}
}

// Parse takes a specification string, and turns parses it into an Argument.
func (p *argumentSpecifcationParser) parse() (parameters.Argument, error) {
	var argument parameters.Argument
	var deep bool

	if tok, _ := p.scan(); tok == LBRACK {
		deep = true
	} else {
		argument.Required = true
		p.unscan()
	}

	if tok, lit := p.scan(); tok == IDENTIFIER {
		argument.Name = strings.ToUpper(lit)
	} else {
		return argument, p.expected("identifier", lit)
	}

	if deep {
		if tok, lit := p.scan(); tok != RBRACK {
			return argument, p.expected("closing bracket", lit)
		}
	}

	if tok, lit := p.scan(); tok != EOF {
		return argument, p.expected("end of spec", lit)
	}

	return argument, nil
}
