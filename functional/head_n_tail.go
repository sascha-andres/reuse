package functional

// Head returns the first element of the provided slice or zero value if the slice is empty.
func Head[T any](rows []T) (t T) {
	if len(rows) == 0 {
		return
	}
	return rows[0]
}

// Tail returns a new slice containing all elements of the input slice except the first one.
// If the input slice is empty, it returns an empty slice.
func Tail[T any](rows []T) (t []T) {
	if len(rows) == 0 {
		return
	}
	return rows[1:]
}
