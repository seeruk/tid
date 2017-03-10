package errhandling

import (
	"errors"
	"fmt"
)

// ErrorStack is an error type that is used to handle multiple errors at once.
type ErrorStack struct {
	errors []error
}

// NewErrorStack creates a new instance of ErrorStack
func NewErrorStack() *ErrorStack {
	return new(ErrorStack)
}

// Add an error to the stack.
func (s *ErrorStack) Add(err error) {
	if err != nil {
		s.errors = append(s.errors, err)
	}
}

// Empty returns true if there are no errors in the stack.
func (s *ErrorStack) Empty() bool {
	return len(s.errors) == 0
}

// Errors returns a summary error if there are any errors in the stack.
func (s *ErrorStack) Errors() error {
	if s.Empty() {
		return nil
	}

	result := "errhandling: The following errors occurred:"

	for _, err := range s.errors {
		result = fmt.Sprintf("%s\n%s", result, err.Error())
	}

	return errors.New(result)
}
