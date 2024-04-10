package functional

import "github.com/sascha-andres/reuse"

// Filter applies the provided function to each element in the input slice and
// returns a new slice containing only the elements for which the function
// returns true. If the function returns an error for any element, Filter
// immediately returns nil and the error.
//
// The input slice `in` holds the original elements.
//
// The function `f` is applied to each element in `in`. It takes an element of
// type `T` as input and returns a boolean value indicating whether the element
// should be included in the output slice, and an error if any occurred.
//
// The returned slice and error value represent the filtered elements and any
// error encountered while performing the filtering operation, respectively.
// If no error occurred, the error value is nil.
//
// Example usage:
//
//	nums := []int{1, 2, 3, 4, 5}
//	evenNums, err := Filter(nums, func(n int) (bool, error) {
//	    if n%2 == 0 {
//	        return true, nil
//	    }
//	    return false, nil
//	})
//	if err != nil {
//	    fmt.Println("Error occurred:", err)
//	}
//	fmt.Println(evenNums) // Output: [2, 4]
//
// This function is generic and works with any type T.
//
// The function signature is:
//
//	func Filter[T any](in []T, f func(T) (bool, error)) ([]T, error)
func Filter[T any](in []T, f func(T) (bool, error)) ([]T, error) {
	var out []T
	var errs []error
	for _, t := range in {
		if r, err := f(t); err == nil && r {
			out = append(out, t)
		} else {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return out, nil
	}
	return out, reuse.MultiError(errs)
}
