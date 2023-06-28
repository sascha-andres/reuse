package reuse

// MultiError is a list of errors
type MultiError []error

// Error returns a string rep of all errors
//
// Hint: probably you want to use errors.Unwrap approach
// this is meant more to collect a number of errors
// on data rows where you do not want (or cannot) return
// on first error
func (m MultiError) Error() string {
	var result string
	for _, err := range m {
		if result == "" {
			result = err.Error()
			continue
		}
		result = result + ";" + err.Error()
	}
	return result
}
