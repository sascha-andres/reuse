package functional

// Map will map one object type to another. This is especially useful in a map reduce
// context
func Map[O, T any](rows []O, mapper func(O) (T, error)) ([]T, error) {
	result := make([]T, 0)
	for _, row := range rows {
		res, err := mapper(row)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}
