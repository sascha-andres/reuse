package reuse

// MultiError is a list of errors
type MultiError []error

// Error returns a string rep of all errors
func (m MultiError) Error() string {
	var result string
	for _, err := range m {
		if result == "" {
			result = err.Error()
		}
		result = result + ";" + err.Error()
	}
	return result
}
