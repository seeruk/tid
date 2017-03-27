package console

import (
	"fmt"
	"io"
)

// Output abstracts application output. This is mainly useful for testing, as a different writer can
// be passed to capture output in an easy to test manner.
type Output struct {
	Writer io.Writer
}

// NewOutput creates a new Output.
func NewOutput(writer io.Writer) *Output {
	return &Output{
		Writer: writer,
	}
}

// Print uses the fmt package's Print with a pre-set writer. Spaces are always added between
// operands. It returns the number of bytes written and any write error encountered.
func (o *Output) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(o.Writer, a...)
}

// Printf uses the fmt package's Printf with a pre-set writer. It returns the number of bytes
// written and any write error encountered.
func (o *Output) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(o.Writer, format, a...)
}

// Println uses the fmt package's Println with a pre-set writer. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes written and any write error
// encountered.
func (o *Output) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(o.Writer, a...)
}
