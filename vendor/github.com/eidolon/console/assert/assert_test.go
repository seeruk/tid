package assert_test

import (
	"errors"
	"testing"

	"github.com/eidolon/console/assert"
)

// mockTester is a simple implementation of the assert.tester interface, to aid in writing tests.
type mockTester struct {
	errorfCalls int
}

// Errorf increments a counter to provide call count information in tests.
func (t *mockTester) Errorf(format string, args ...interface{}) {
	t.errorfCalls = t.errorfCalls + 1
}

func TestTrue(t *testing.T) {
	t.Run("should not end testing if the condition is truthy", func(t *testing.T) {
		tester := new(mockTester)

		assert.True(tester, true, "")

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the condition is falsey", func(t *testing.T) {
		tester := new(mockTester)

		assert.True(tester, false, "")

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}

func TestFalse(t *testing.T) {
	t.Run("should not end testing if the condition is falsey", func(t *testing.T) {
		tester := new(mockTester)

		assert.False(tester, false, "")

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the condition is truthy", func(t *testing.T) {
		tester := new(mockTester)

		assert.False(tester, true, "")

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}

func TestEqual(t *testing.T) {
	t.Run("should not end testing if the expected equals the actual", func(t *testing.T) {
		tester := new(mockTester)

		assert.Equal(tester, 42, 42)

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the expected does not equal the actual", func(t *testing.T) {
		tester := new(mockTester)

		assert.Equal(tester, 24, 42)

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("should not end testing if the expected does not equal the actual", func(t *testing.T) {
		tester := new(mockTester)

		assert.NotEqual(tester, 24, 42)

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the expected does equal the actual", func(t *testing.T) {
		tester := new(mockTester)

		assert.NotEqual(tester, 42, 42)

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}

func TestOK(t *testing.T) {
	t.Run("should not end testing if the given error is nil", func(t *testing.T) {
		tester := new(mockTester)

		assert.OK(tester, nil)

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the given error is not nil", func(t *testing.T) {
		tester := new(mockTester)

		assert.OK(tester, errors.New("Testing"))

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}

func TestNotOK(t *testing.T) {
	t.Run("should not end testing if the given error is not nil", func(t *testing.T) {
		tester := new(mockTester)

		assert.NotOK(tester, errors.New("Testing"))

		if tester.errorfCalls != 0 {
			t.Errorf(
				"Expected Errorf not to have been called, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})

	t.Run("should end testing if the given error is nil", func(t *testing.T) {
		tester := new(mockTester)

		assert.NotOK(tester, nil)

		if tester.errorfCalls != 1 {
			t.Errorf(
				"Expected Errorf to have been called once, it had been called %d times.",
				tester.errorfCalls,
			)
		}
	})
}
