package functional

import (
	"reflect"
	"testing"
)

var filterTests = []struct {
	name string
	in   []int
	out  []int
	f    func(int) (bool, error)
}{
	{
		name: "is even",
		in:   []int{1, 2, 3, 4, 5},
		out:  []int{2, 4},
		f: func(n int) (bool, error) {
			if n%2 == 0 {
				return true, nil
			}
			return false, nil
		},
	},
}

// TestFilter is a test function that iterates over a list of filter tests
// and runs each test case using the Filter function. It checks if the result
// matches the expected output and if any error occurred during the filtering
// process. This function is used to verify the correctness of the Filter function.
//
// The input slice `filterTests` contains a list of test cases defined as a struct,
// which consists of the test name, the input slice, the expected output slice,
// and the filter function. Each test case applies the filter function to the input
// slice and compares the result with the expected output.
//
// To run the tests, TestFilter invokes the t.Run function with the test name as
// subtests, providing isolation for each test case. Inside each subtest, it calls
// the Filter function with the input and filter function from the current test case,
// and checks if the returned result matches the expected output. If an error occurs
// during the filtering process, it reports an unexpected error. If the result does
// not match the expected output, it reports an error indicating the expected and
// actual results.
//
// Example usage:
//
//	go test -run TestFilter
//
// This function is used for testing the generic Filter function and does not provide
// any public API.
//
// To understand how the Filter function works, please refer to the documentation
// of the Filter function.
func TestFilter(t *testing.T) {
	for _, test := range filterTests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Filter(test.in, test.f)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, test.out) {
				t.Errorf("expected %v, got %v", test.out, result)
			}
		})
	}
}
