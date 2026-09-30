# lo

`slog`-logging variants of the "must" pattern, for code that wants to
log-and-continue rather than propagate every error.

```go
logger := slog.Default()

value := lo.MustGet(logger, "reading config", func() (Config, error) {
    return loadConfig()
})
// on error: logged via logger.Error, value is the zero Config

lo.MustRun(logger, "flushing cache", func() error {
    return cache.Flush()
})
// on error: logged, execution continues

value := lo.Must(loadConfig()) // panics if err != nil
```

See also [te](../te/) for the equivalent pattern against `*testing.T`
instead of `*slog.Logger`.
