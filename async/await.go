package async

import (
	"fmt"
	"runtime/debug"
	"context"
)

var ErrCancelled = fmt.Errorf("cancelled")

// Future represents a value that will be available at some point in the future
type Future[T any] struct {
	// await returns the value of the Future
	await func() (T, error)
}

// Await waits for the Future to complete and returns the result
func (f *Future[T]) Await() (T, error) {
	return f.await()
}

// Async wraps a function returning a value of type T and returns a Future[T]
func Async[T any](ctx context.Context, f func(context.Context) (T, error)) *Future[T] {
	var result T
	var err error

	done := make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					err = fmt.Errorf("recovering from error: %w\n%s", x, debug.Stack())
				default:
					err = fmt.Errorf("panic: %v\n%s", x, debug.Stack())
				}
			}
		}()

		result, err = f(ctx)
	}()

	return &Future[T]{
		await: func() (T, error) {
			<-done
			return result, err
		},
	}
}
