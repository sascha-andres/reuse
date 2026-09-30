# te

`testing`-failing variants of the "must" pattern, for test helpers that
call functions returning `(T, error)` / `error`.

```go
func TestSomething(t *testing.T) {
    cfg := te.MustGet(t, "loading fixture", func() (Config, error) {
        return loadConfig("testdata/config.json")
    })
    // t.Fatalf on error, stops the test immediately

    result := te.ShouldGet(t, "computing result", func() (int, error) {
        return compute(cfg)
    })
    // t.Error + t.Fail on error, test continues

    te.MustRun(t, "writing output", func() error {
        return writeOutput(result)
    })
    te.ShouldRun(t, "cleaning up", func() error {
        return cleanup()
    })
}
```

`Must*` stops the test (`t.Fatalf`) on error; `Should*` records a failure
(`t.Error` + `t.Fail`) and lets the test continue. See also [lo](../lo/)
for the equivalent pattern against `*slog.Logger` instead of
`*testing.T`.
