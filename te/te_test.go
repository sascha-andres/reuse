package te

import (
	"fmt"
	"testing"
)

// ExampleMustGet demonstrates the usage of MustGet with functions returning values and errors.
// It shows how a successful result is returned and how an error causes the test to fail immediately.
func ExampleMustGet() {
	b := MustGet(&testing.T{}, "test", func() (string, error) { return "a", nil })
	// Output: b == a
	_ = MustGet(&testing.T{}, "test", func() (bool, error) { return b == "b", fmt.Errorf("error") })
	// Output: Failed test
}

// ExampleGetOrFailNow demonstrates the usage of MustGet with functions returning values and errors.
// It shows how a successful result is returned and how an error causes the test to fail immediately.
func ExampleShouldGet() {
	b := ShouldGet(&testing.T{}, "test", func() (string, error) { return "a", nil })
	// Output: b == a
	_ = ShouldGet(&testing.T{}, "test", func() (bool, error) { return b == "b", fmt.Errorf("error") })
	// Output: Failed test
}

// ExampleMustRun demonstrates the usage of the MustRun function to handle error-prone test logic and report failures appropriately.
func ExampleMustRun() {
	MustRun(&testing.T{}, "test", func() error { return nil })
	MustRun(&testing.T{}, "test", func() error { return fmt.Errorf("error") })
	// Output: Failed test
}

// ExampleShouldRun demonstrates the usage of the ShouldRun function, checking for error handling and test failure reporting.
func ExampleShouldRun() {
	ShouldRun(&testing.T{}, "test", func() error { return nil })
	ShouldRun(&testing.T{}, "test", func() error { return fmt.Errorf("error") })
	// Output: Failed test
}
