package specification_test

import (
	"testing"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
	"github.com/eidolon/console/specification"
)

func TestParseOptionSpecification(t *testing.T) {
	t.Run("should set the name(s)", func(t *testing.T) {
		option, err := specification.ParseOptionSpecification("--galaxy-quest")
		assert.OK(t, err)
		assert.True(t, len(option.Names) == 1, "Expected an option name")

		name := option.Names[0]
		assert.Equal(t, "galaxy-quest", name)

		option, err = specification.ParseOptionSpecification("-g, --galaxy-quest")
		assert.OK(t, err)
		assert.True(t, len(option.Names) == 2, "Expected 2 option names")

		shortName := option.Names[0]
		longName := option.Names[1]
		assert.Equal(t, "g", shortName)
		assert.Equal(t, "galaxy-quest", longName)

		option, err = specification.ParseOptionSpecification("-a, -b, -c, --dee, --ee, --eff")
		assert.OK(t, err)
		assert.True(t, len(option.Names) == 6, "Expected 6 option names")
	})

	t.Run("should allow no value mode to be set", func(t *testing.T) {
		option, err := specification.ParseOptionSpecification("--galaxy-quest")
		assert.OK(t, err)
		assert.Equal(t, parameters.OptionValueNone, option.ValueMode)
	})

	t.Run("should allow option values", func(t *testing.T) {
		option, err := specification.ParseOptionSpecification("--galaxy-quest-2[=ALAN_RICKMAN]")
		assert.OK(t, err)
		assert.Equal(t, parameters.OptionValueOptional, option.ValueMode)
	})

	t.Run("should allow required values", func(t *testing.T) {
		option, err := specification.ParseOptionSpecification("--galaxy-quest=ALAN_RICKMAN")
		assert.OK(t, err)
		assert.Equal(t, parameters.OptionValueRequired, option.ValueMode)
	})

	t.Run("should error when given an invalid long option name", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("--$$$")
		assert.NotOK(t, err)
	})

	t.Run("should error when given an invalid short option name", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("-$")
		assert.NotOK(t, err)
	})

	t.Run("should error when given an short option name that's too long", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("-abc")
		assert.NotOK(t, err)
	})

	t.Run("should error when given an invalid specification", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("abc")
		assert.NotOK(t, err)
	})

	t.Run("should error if no identifier is given for the value name", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("--galaxy-quest=")
		assert.NotOK(t, err)

		_, err = specification.ParseOptionSpecification("--galaxy-quest[=]")
		assert.NotOK(t, err)
	})

	t.Run("should expect EQUALS if an LBRACK is found", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("--galaxy-quest[]")
		assert.NotOK(t, err)

		_, err = specification.ParseOptionSpecification("--galaxy-quest[")
		assert.NotOK(t, err)
	})

	t.Run("should expect an RBRACK if an LBRACK is given", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("--galaxy-quest-2[=ALAN_RICKMAN")
		assert.NotOK(t, err)
	})

	t.Run("should expect EOF after value identifier is no LBRACK is given", func(t *testing.T) {
		_, err := specification.ParseOptionSpecification("--galaxy-quest-2=ALAN_RICKMAN]")
		assert.NotOK(t, err)
	})
}
