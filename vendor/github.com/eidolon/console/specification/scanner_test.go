package specification_test

import (
	"strings"
	"testing"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/specification"
)

func TestNewScanner(t *testing.T) {
	scanner := specification.NewScanner(strings.NewReader(""))
	assert.True(t, scanner != nil, "Scanner should not be nil")
}

func TestScanner(t *testing.T) {
	createScanner := func(input string) *specification.Scanner {
		return specification.NewScanner(strings.NewReader(input))
	}

	t.Run("Scan()", func(t *testing.T) {
		t.Run("should handle end of input", func(t *testing.T) {
			scanner := createScanner("=")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.EQUALS, tok)
			assert.Equal(t, "=", val)

			tok, val = scanner.Scan()
			assert.Equal(t, specification.EOF, tok)
			assert.Equal(t, "", val)
		})

		t.Run("should be able to scan left-brackets ([)", func(t *testing.T) {
			scanner := createScanner("[")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.LBRACK, tok)
			assert.Equal(t, "[", val)
		})

		t.Run("should be able to scan right-brackets (])", func(t *testing.T) {
			scanner := createScanner("]")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.RBRACK, tok)
			assert.Equal(t, "]", val)
		})

		t.Run("should be able to scan commas (,)", func(t *testing.T) {
			scanner := createScanner(",")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.COMMA, tok)
			assert.Equal(t, ",", val)
		})

		t.Run("should be able to scan equals (=)", func(t *testing.T) {
			scanner := createScanner("=")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.EQUALS, tok)
			assert.Equal(t, "=", val)
		})

		t.Run("should be able to scan hyphens (-)", func(t *testing.T) {
			scanner := createScanner("-")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.HYPHEN, tok)
			assert.Equal(t, "-", val)
		})

		t.Run("should be able to scan whitespace", func(t *testing.T) {
			scanner := createScanner(" 	")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.WS, tok)
			assert.Equal(t, " 	", val)
		})

		t.Run("should only scan up to available whitespace characters", func(t *testing.T) {
			scanner := createScanner("      hello")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.WS, tok)
			assert.Equal(t, "      ", val)
		})

		t.Run("should be able to scan identifiers", func(t *testing.T) {
			scanner := createScanner("I_am_an_identifier-scanner")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.IDENTIFIER, tok)
			assert.Equal(t, "I_am_an_identifier-scanner", val)
		})

		t.Run("should only scan up to available identifier characters", func(t *testing.T) {
			scanner := createScanner("SOME_IDENTIFIER=")

			tok, val := scanner.Scan()
			assert.Equal(t, specification.IDENTIFIER, tok)
			assert.Equal(t, "SOME_IDENTIFIER", val)

			tok, val = scanner.Scan()
			assert.Equal(t, specification.EQUALS, tok)
			assert.Equal(t, "=", val)
		})

		t.Run("should handle illegal characters", func(t *testing.T) {
			input := "£$%"

			scanner := createScanner(input)

			tok, val := scanner.Scan()
			assert.Equal(t, specification.ILLEGAL, tok)
			assert.Equal(t, "£", val)

			tok, val = scanner.Scan()
			assert.Equal(t, specification.ILLEGAL, tok)
			assert.Equal(t, "$", val)

			tok, val = scanner.Scan()
			assert.Equal(t, specification.ILLEGAL, tok)
			assert.Equal(t, "%", val)
		})
	})
}
