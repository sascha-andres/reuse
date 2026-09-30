# circuitbreaker

A minimal circuit breaker around a `func(id int) error` work function.

## States

`Closed` (calls run normally) → `Open` (calls are rejected without
running, once `maxFailures` consecutive failures are reached) → after
`timeout` has elapsed, `Half-Open` (one call is allowed through to probe
recovery) → back to `Closed` on success, back to `Open` on failure.

## Usage

`Call` is designed to be run as a goroutine, reporting its result on a
channel and signaling a `sync.WaitGroup`:

```go
cb := circuitbreaker.NewCircuitBreaker(3, 10*time.Second, func(id int) error {
    return doWork(id)
})

var wg sync.WaitGroup
taskDone := make(chan circuitbreaker.Task)

for id := 0; id < n; id++ {
    wg.Add(1)
    go cb.Call(&wg, taskDone, id)
}

go func() {
    wg.Wait()
    close(taskDone)
}()

for t := range taskDone {
    if t.Status {
        fmt.Println("task", t.Id, "succeeded")
    } else {
        fmt.Println("task", t.Id, "failed or rejected")
    }
}
```

A `Task` with `Status == false` covers both an actual failure of `fn` and
a call rejected outright because the breaker was `Open`; the two aren't
distinguished on the channel.
