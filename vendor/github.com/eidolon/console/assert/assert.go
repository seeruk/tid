package assert

import (
	"os"
	"reflect"
	"runtime"
	"strings"
)

// tester provides a simpler testing interface for use in mocks.
type tester interface {
	Errorf(format string, args ...interface{})
}

// True asserts that the given condition is truthy.
func True(t tester, condition bool, message string) {
	_, file, line, _ := runtime.Caller(1)

	if !condition {
		t.Errorf(
			"assert: %s:%d: %s",
			trimLocation(file),
			line,
			message,
		)
	}
}

// False asserts that the given condition is falsey.
func False(t tester, condition bool, message string) {
	_, file, line, _ := runtime.Caller(1)

	if condition {
		t.Errorf(
			"assert: %s:%d: %s",
			trimLocation(file),
			line,
			message,
		)
	}
}

// Equal asserts that the given expected and actual arguments are equal.
func Equal(t tester, expected, actual interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"assert: %s:%d: Expected %v (type %v) to equal %v (type %v).",
			trimLocation(file),
			line,
			actual,
			reflect.TypeOf(actual),
			expected,
			reflect.TypeOf(expected),
		)
	}
}

// NotEqual asserts that the given expected and actual arguments are not equal.
func NotEqual(t tester, expected, actual interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"assert: %s:%d: Expected %v (type %v) not to equal got %v (type %v).",
			trimLocation(file),
			line,
			actual,
			reflect.TypeOf(actual),
			expected,
			reflect.TypeOf(expected),
		)
	}
}

// OK checks that an error was not produced, and reacts accordingly.
func OK(t tester, err error) {
	_, file, line, _ := runtime.Caller(1)

	if err != nil {
		t.Errorf(
			"assert: %s:%d: Unexpected error: '%s'.",
			trimLocation(file),
			line,
			err.Error(),
		)
	}
}

// NotOK checks that an error was produced, and reacts accordingly.
func NotOK(t tester, err error) {
	_, file, line, _ := runtime.Caller(1)

	if err == nil {
		t.Errorf(
			"assert: %s:%d: Expected error, but none was given.",
			trimLocation(file),
			line,
		)
	}
}

// trimLocation takes an absolute file path, and returns a much shorter relative path.
func trimLocation(file string) string {
	// Alright, I know, this error is ignored, I'm a terrible person, yada yada yada. It'll just
	// panic anyway, and I don't think that this is testable otherwise...
	cwd, _ := os.Getwd()

	return strings.Replace(file, cwd+"/", "", -1)
}
