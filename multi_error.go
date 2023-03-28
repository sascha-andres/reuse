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

// Add will add an error to the list
func (m MultiError) Add(err error) MultiError {
	return append(m, err)
}
