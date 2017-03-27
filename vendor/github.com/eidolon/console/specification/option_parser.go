package specification

import (
	"strings"

	"github.com/eidolon/console/parameters"
)

// ParseOptionSpecification parses an option spec string and produces an Option.
func ParseOptionSpecification(spec string) (parameters.Option, error) {
	scanner := NewScanner(strings.NewReader(spec))
	parser := newOptionSpecifcationParser(scanner)

	return parser.parse()
}

// optionSpecificationParser parses option specification strings.
type optionSpecificationParser struct {
	parser
}

// newOptionSpecifcationParser creates a new option specification parser.
func newOptionSpecifcationParser(scanner *Scanner) *optionSpecificationParser {
	parser := parser{
		scanner: scanner,
	}

	return &optionSpecificationParser{
		parser: parser,
	}
}

// parse takes a spec string, and turns parses it into an Option.
func (p *optionSpecificationParser) parse() (parameters.Option, error) {
	var option parameters.Option

	for {
		tok, lit := p.scan()

		if tok == HYPHEN {
			name, err := p.parseOptionName()
			if err != nil {
				return option, err
			}

			option.Names = append(option.Names, name)
		} else {
			return option, p.expected("hyphen", lit)
		}

		if tok, _ := p.scan(); tok != COMMA {
			p.unscan()
			break
		}
	}

	tok, _ := p.scan()

	if tok == LBRACK || tok == EQUALS {
		var deep bool

		// Handle opening bracket
		if tok == LBRACK {
			option.ValueMode = parameters.OptionValueOptional
			deep = true
		} else {
			option.ValueMode = parameters.OptionValueRequired
			p.unscan()
		}

		if tok, lit := p.scan(); tok == EQUALS {
			if tok, lit := p.scan(); tok == IDENTIFIER {
				option.ValueName = strings.ToUpper(lit)
			} else {
				return option, p.expected("identifier", lit)
			}
		} else {
			return option, p.expected("equals", lit)
		}

		if deep {
			if tok, lit := p.scan(); tok != RBRACK {
				return option, p.expected("closing bracket", lit)
			}
		}
	} else {
		option.ValueMode = parameters.OptionValueNone
	}

	if tok, lit := p.scan(); tok != EOF {
		return option, p.expected("end of spec", lit)
	}

	return option, nil
}

// parseOptionName attempts to parse a set of tokens that form an option name.
func (p *optionSpecificationParser) parseOptionName() (string, error) {
	tok, lit := p.scan()

	if tok == HYPHEN {
		return p.parseLongOptionName()
	} else if tok == IDENTIFIER {
		return p.parseShortOptionName()
	} else {
		return lit, p.expected("hyphen or identifier", lit)
	}
}

// parseLongOptionName attempts to parse a set of tokens that form a long option name.
func (p *optionSpecificationParser) parseLongOptionName() (string, error) {
	tok, lit := p.scan()

	if tok == IDENTIFIER {
		return lit, nil
	}

	return lit, p.expected("identifier", lit)
}

// parseShortOptionName attempts to parse a set of tokens that form a short option name.
func (p *optionSpecificationParser) parseShortOptionName() (string, error) {
	p.unscan()

	_, lit := p.scan()

	if len(lit) == 1 {
		return lit, nil
	}

	return lit, p.expectedLen("short option identifier", 1, len(lit))
}
