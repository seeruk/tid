package specification

import (
	"bytes"
	"strings"
	"unicode"
)

// Token identifies a type of lexical token.
type Token int

// eof represents the end of input.
const eof = rune(0)

// The possible token types.
const (
	ILLEGAL Token = iota
	EOF

	LBRACK     // [
	RBRACK     // ]
	COMMA      // ,
	EQUALS     // =
	HYPHEN     // -
	WS         // Whitespace
	IDENTIFIER // A-z_-
)

// Scanner is a lexical scanner for parameter specifications.
type Scanner struct {
	reader *strings.Reader
}

// NewScanner creates a new specification scanner.
func NewScanner(reader *strings.Reader) *Scanner {
	return &Scanner{
		reader: reader,
	}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (Token, string) {
	r := s.read()

	switch r {
	case eof:
		return EOF, ""
	case '[':
		return LBRACK, string(r)
	case ']':
		return RBRACK, string(r)
	case ',':
		return COMMA, string(r)
	case '=':
		return EQUALS, string(r)
	case '-':
		return HYPHEN, string(r)
	}

	if isWhitespace(r) {
		s.unread()

		return s.scanWhitespace()
	} else if isAlphaNumeric(r) {
		s.unread()

		return s.scanIdentifier()
	}

	return ILLEGAL, string(r)
}

// scanIdentifier consumes the current rune and all contiguous identifier runes.
func (s *Scanner) scanIdentifier() (Token, string) {
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isAlphaNumeric(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return IDENTIFIER, buf.String()
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer

	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// read reads the next rune from the buffered reader. Returns the rune(0) if an error occurs (or if
// io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()

	if err != nil {
		return eof
	}

	return ch
}

// unread rewinds the reader to re-read the previously read rune again.
func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

// isAlphaNumeric reports whether the given rune is an alphanumeric character (or - / _).
func isAlphaNumeric(r rune) bool {
	return r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isWhitespace reports whether the given rune is a whitespace character.
func isWhitespace(r rune) bool {
	return r == '	' || r == ' ' || r == '\n' || r == '\r' || r == '\t' || unicode.IsSpace(r)
}
