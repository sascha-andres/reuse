package functional

import "sync"

// ValueCell wraps the interaction functions into a struct so that CreateValueCell
// may return a single value only
type ValueCell[T comparable] struct {
	// Get returns the current value
	Get func() T
	// Set updates the curent value
	Set func(T)
	// AddWatcher allows adding a watcher to get notified on value change// AddWatcher allows adding a watcher to get notified on value change
	AddWatcher func(func(T, T))
}

// CreateValueCell creates a ValueCell. A value cells simply wrap a variable with
// two simple operations (Get and Set). Access is locked with a sync.Mutex.
// To get notified about value changes use AddWatcher to provide watchers
//
// Usage:
//
//	c1 := CreateValueCell(1)
//	c1.AddWatcher(func(oldValue, newValue int) {
//	  fmt.Printf("\n ** %d -> %d **\n", oldValue, newValue)
//	})
//	c2 := CreateValueCell(2)
//	fmt.Printf("%d\n", c1.Get()+c2.Get())
//	c1.Set(2)
//	fmt.Printf("%d", c1.Get()+c2.Get())
func CreateValueCell[T comparable](initial T) ValueCell[T] {
	value := initial
	watchers := make([]func(oldValue, newValue T), 0)
	var mu sync.Mutex
	g := func() T {
		mu.Lock()
		defer mu.Unlock()
		return value
	}
	u := func(newValue T) {
		mu.Lock()
		defer mu.Unlock()
		oldValue := value
		if oldValue != newValue {
			value = newValue
			for _, watcher := range watchers {
				watcher(oldValue, newValue)
			}
		}
	}
	w := func(watcher func(oldValue, newValue T)) {
		watchers = append(watchers, watcher)
	}
	c := ValueCell[T]{
		Get:        g,
		Set:        u,
		AddWatcher: w,
	}
	return c
}
