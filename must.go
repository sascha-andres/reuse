package reuse

// Must panics if err is not nil
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
