# functional

Generic slice/map helpers and two small reactive-cell primitives.

## Slice/map helpers

```go
evens, err := functional.Filter([]int{1, 2, 3, 4}, func(n int) (bool, error) {
    return n%2 == 0, nil
})
// evens == []int{2, 4}

strs, err := functional.Map([]int{1, 2, 3}, func(n int) (string, error) {
    return strconv.Itoa(n), nil
})
// strs == []string{"1", "2", "3"}

first := functional.Head([]int{1, 2, 3}) // 1
rest := functional.Tail([]int{1, 2, 3})  // []int{2, 3}
```

`Filter` collects every error encountered (rather than stopping at the
first) and returns them together as a `reuse.MultiError`; elements whose
callback errored are simply excluded from the result. `Map` stops and
returns immediately on the first error.

`GroupBy` groups a slice of `map[string]T` rows by the values of one or
more keys, joined with `_`; `GroupByFunc` groups an arbitrary slice by a
key function:

```go
groups, err := functional.GroupBy(rows, "country", "city")
// err is ErrMissingColumn if a row is missing one of the given keys

byLen, err := functional.GroupByFunc(words, func(w string) (int, error) {
    return len(w), nil
})
```

## Value/formula cells

`ValueCell` wraps a mutex-guarded value with `Get`/`Set`/`AddWatcher`;
`Set` returns `false` instead of blocking if the lock is currently held.
`FormulaCell` derives from a `ValueCell` (or another `FormulaCell`'s
watcher) and recomputes whenever its upstream changes:

```go
c1 := functional.CreateValueCell(1)
doubled := functional.CreateFormulaCell(c1, func(old, new int) int {
    return new * 2
})
doubled.AddWatcher(func(old, new int) {
    fmt.Println(old, "->", new)
})
c1.Set(5) // triggers the watcher with the recomputed value
```

A `FormulaCell` only watches a single upstream `ValueCell`; `calculation`
receives that upstream's old and new value on each change.
