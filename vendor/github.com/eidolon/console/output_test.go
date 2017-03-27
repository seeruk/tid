package console_test

import (
	"bytes"
	"testing"

	"github.com/eidolon/console"
	"github.com/eidolon/console/assert"
)

func TestNewOutput(t *testing.T) {
	output := console.NewOutput(&bytes.Buffer{})

	assert.True(t, output != nil, "Output should not be nil")
}

func TestOutput(t *testing.T) {
	// Realistically, these tests are a little lacking, but that's mainly because we're pretty much
	// transparently passing everything over to the fmt package. "So why are you bothering to do
	// that?", simply so we can customise the writer in the application, and not have to specify it
	// every time we write.

	t.Run("Print()", func(t *testing.T) {
		t.Run("should handle printing", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			message := "Hello, World!"

			nbytes, err := output.Print(message)

			assert.OK(t, err)
			assert.Equal(t, 13, nbytes)
			assert.Equal(t, message, buffer.String())
		})

		t.Run("should handle printing multiple values", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			message1 := "Go "
			message2 := "testing!"

			nbytes, err := output.Print(message1, message2)

			assert.OK(t, err)
			assert.Equal(t, 11, nbytes)
			assert.Equal(t, message1+message2, buffer.String())
		})
	})

	t.Run("Printf()", func(t *testing.T) {
		t.Run("should handle printing", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			nbytes, err := output.Printf("Hello, %s!", "World")

			assert.OK(t, err)
			assert.Equal(t, 13, nbytes)
			assert.Equal(t, "Hello, World!", buffer.String())
		})

		t.Run("should handle multiple parameters", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			message := "Hello %s, the answer to life the universe and everything is %d!"
			param1 := "Go"
			param2 := 42

			expected := "Hello Go, the answer to life the universe and everything is 42!"

			nbytes, err := output.Printf(message, param1, param2)

			assert.OK(t, err)
			assert.Equal(t, 63, nbytes)
			assert.Equal(t, expected, buffer.String())
		})
	})

	t.Run("Println()", func(t *testing.T) {
		t.Run("should handle printing", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			message := "Hello, World!"
			expected := "Hello, World!\n"

			nbytes, err := output.Println(message)

			assert.OK(t, err)
			assert.Equal(t, 14, nbytes)
			assert.Equal(t, expected, buffer.String())
		})

		t.Run("should handle printing multiple values", func(t *testing.T) {
			buffer := bytes.Buffer{}
			output := console.NewOutput(&buffer)

			message1 := "Go "
			message2 := "testing!"

			// Remember, fmt.Println adds spaces between things.
			expected := "Go  testing!\n"

			nbytes, err := output.Println(message1, message2)

			assert.OK(t, err)
			assert.Equal(t, 13, nbytes)
			assert.Equal(t, expected, buffer.String())
		})
	})
}
