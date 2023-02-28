package functional

// FormulaCell wraps the interaction functions into a struct so that CreateFormulaCell
// may return a single value only
type FormulaCell[T comparable] struct {
	// Get returns the current value
	Get func() T
	// AddWatcher allows adding a watcher to get notified on value change// AddWatcher allows adding a watcher to get notified on value change
	AddWatcher func(func(T, T))
}

// CreateFormulaCell creates a FormulaCell. A formula cell depends on a ValueCell
// and operates on it. When ValueCell's value changes (FormulaCell adds a watcher)
// it will receive old and new value and can calculate on it. Initially calculation
// will be called with the initial value of upstream as old and new
//
//	c1 := CreateValueCell(1)
//	f := CreateFormulaCell(c1, func(op1, op2 int) int {
//		return op1 + op2
//	})
//	f.AddWatcher(func(oldValue, newValue int) {
//		fmt.Printf("\n ++ %d -> %d ++\n", oldValue, newValue)
//	})
//	c1.Set(2)
func CreateFormulaCell[T comparable](upstream ValueCell[T], calculation func(T, T) T) FormulaCell[T] {
	value := CreateValueCell(calculation(upstream.Get(), upstream.Get()))
	watchers := make([]func(oldValue, newValue T), 0)
	upstream.AddWatcher(func(oldValue T, newValue T) {
		calculatedOld := value.Get()
		value.Set(calculation(oldValue, newValue))
		if value.Get() != calculatedOld {
			for _, watcher := range watchers {
				watcher(calculatedOld, value.Get())
			}
		}
	})
	g := value.Get
	w := func(watcher func(oldValue, newValue T)) {
		watchers = append(watchers, watcher)
	}
	f := FormulaCell[T]{
		Get:        g,
		AddWatcher: w,
	}
	return f
}
