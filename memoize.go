package reuse

// Memoize returns access to a cache
//
// parameters
// fn - function to query for value
//
// return values:
// 1 - value, in case of error undefined (depending on fn)
// 2 - indicating cached (true) or miss (false)
// 3 - error returned from fn
func Memoize[T comparable](fn func(T) (T, error)) func(T) (T, bool, error) {
	cache := make(map[T]T)
	return func(n T) (T, bool, error) {
		if v, ok := cache[n]; ok {
			return v, true, nil
		}
		val, err := fn(n)
		if err != nil {
			return val, false, err
		}
		cache[n] = val
		return cache[n], false, nil
	}
}
