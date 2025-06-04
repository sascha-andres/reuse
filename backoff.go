package reuse

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

// delays defines a slice of durations used for exponential backoff in retry logic, with fixed and increasing intervals.
var delays = []time.Duration{
	1 * time.Second, 2 * time.Second,
	4 * time.Second, 8 * time.Second,
	16 * time.Second, 32 * time.Second,
	60 * time.Second, 60 * time.Second,
	60 * time.Second, 60 * time.Second,
}

// SetDelays sets the custom delay durations for retry logic by replacing the default delay slice.
func SetDelays(d []time.Duration) {
	delays = d
}

// RequestCtx defines a function that processes a request with a context and returns an error.
type RequestCtx func(ctx context.Context) error

// Request defines a function that processes a request and returns an error.
type Request func() error

// DoCtx executes a request until it succeeds or the context is canceled.
// It will retry the request with exponential backoff.
func DoCtx(ctx context.Context, requestCtx RequestCtx) error {
	for _, delay := range delays {
		err := requestCtx(ctx)
		if err == nil {
			return nil
		}

		delay = multiplyDuration(delay, 0.75+rand.Float64()*0.5) // ±25%
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("failed after %d attempts", len(delays))
}

// Do executes a request until it succeeds or the context is canceled.
// It will retry the request with exponential backoff.
func Do(ctx context.Context, request Request) error {
	for _, delay := range delays {
		err := request()
		if err == nil {
			return nil
		}

		delay = multiplyDuration(delay, 0.75+rand.Float64()*0.5) // ±25%
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("failed after %d attempts", len(delays))
}

// multiplyDuration scales a time.Duration `d` by a float64 multiplier `mul` and returns the resulting duration.
func multiplyDuration(d time.Duration, mul float64) time.Duration {
	return time.Duration(float64(d) * mul)
}
