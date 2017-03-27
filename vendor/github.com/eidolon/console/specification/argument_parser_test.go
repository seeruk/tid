package specification_test

import (
	"testing"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/specification"
)

func TestParseArgumentSpecification(t *testing.T) {
	t.Run("should set the name", func(t *testing.T) {
		argument, err := specification.ParseArgumentSpecification("GALAXY_QUEST")
		assert.OK(t, err)
		assert.Equal(t, "GALAXY_QUEST", argument.Name)
	})

	t.Run("should set whether or not the argument is required", func(t *testing.T) {
		argument, err := specification.ParseArgumentSpecification("GALAXY_QUEST")
		assert.OK(t, err)
		assert.Equal(t, true, argument.Required)

		argument, err = specification.ParseArgumentSpecification("[MEMENTO]")
		assert.OK(t, err)
		assert.Equal(t, false, argument.Required)
	})

	t.Run("should expect a close bracket if an opening one is given", func(t *testing.T) {
		_, err := specification.ParseArgumentSpecification("[MEMENTO")
		assert.NotOK(t, err)
	})

	t.Run("should not expect any whitespace", func(t *testing.T) {
		_, err := specification.ParseArgumentSpecification("GALAXY QUEST")
		assert.NotOK(t, err)
	})

	t.Run("should always expect an identifier", func(t *testing.T) {
		_, err := specification.ParseArgumentSpecification("")
		assert.NotOK(t, err)

		_, err = specification.ParseArgumentSpecification("[]")
		assert.NotOK(t, err)
	})
}
