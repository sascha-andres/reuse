package te

import (
	"fmt"
	"testing"
)

// MustGet executes a function returning a value and error, failing the test immediately on error with a given message.
func MustGet[T any](t *testing.T, msg string, f func() (T, error)) T {
	r, err := f()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
	return r
}

// ShouldGet executes a function returning a error, failing the test on error with a given message.
func ShouldGet[T any](t *testing.T, msg string, f func() (T, error)) T {
	r, err := f()
	if err != nil {
		t.Error(fmt.Sprintf("%s: ", msg), err)
		t.Fail()
	}
	return r
}

// MustRun executes the provided function and reports a test failure with a message if the function returns an error
// and fails immediately.
func MustRun(t *testing.T, msg string, f func() error) {
	err := f()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// ShouldRun executes the provided function and reports a test failure with a message if the function returns an error.
func ShouldRun(t *testing.T, msg string, f func() error) {
	err := f()
	if err != nil {
		t.Error(fmt.Sprintf("%s: ", msg), err)
		t.Fail()
	}
}
