# async

A Go package that provides async/await-like functionality using generics.

## Overview

The `async` package allows you to execute functions asynchronously and wait for their results in a type-safe manner, similar to async/await patterns found in other languages. It uses Go generics to provide compile-time type safety for asynchronous operations.

## Features

- Type-safe asynchronous execution using Go generics
- Automatic panic recovery with stack traces
- Multiple awaits on the same Future
- Context support for cancellation and value propagation
- Works with any type: primitives, structs, pointers, slices, and maps

## Installation

```bash
go get github.com/andres/reuse/async
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/andres/reuse/async"
)

func main() {
    ctx := context.Background()

    // Create a Future by wrapping an async function
    future := async.Async(ctx, func(ctx context.Context) (int, error) {
        // Simulate some work
        return 42, nil
    })

    // Wait for the result
    result, err := future.Await()
    if err != nil {
        panic(err)
    }

    fmt.Println(result) // Output: 42
}
```

### Working with Different Types

```go
// String
future := async.Async(ctx, func(ctx context.Context) (string, error) {
    return "hello world", nil
})

// Struct
type Person struct {
    Name string
    Age  int
}

future := async.Async(ctx, func(ctx context.Context) (Person, error) {
    return Person{Name: "Alice", Age: 30}, nil
})

// Pointer
future := async.Async(ctx, func(ctx context.Context) (*Data, error) {
    return &Data{Value: 100}, nil
})

// Slice
future := async.Async(ctx, func(ctx context.Context) ([]int, error) {
    return []int{1, 2, 3, 4, 5}, nil
})

// Map
future := async.Async(ctx, func(ctx context.Context) (map[string]int, error) {
    return map[string]int{"one": 1, "two": 2}, nil
})
```

### Concurrent Execution

```go
// Execute multiple operations concurrently
future1 := async.Async(ctx, func(ctx context.Context) (int, error) {
    time.Sleep(100 * time.Millisecond)
    return 1, nil
})

future2 := async.Async(ctx, func(ctx context.Context) (int, error) {
    time.Sleep(100 * time.Millisecond)
    return 2, nil
})

future3 := async.Async(ctx, func(ctx context.Context) (int, error) {
    time.Sleep(100 * time.Millisecond)
    return 3, nil
})

// Wait for all results
result1, _ := future1.Await()
result2, _ := future2.Await()
result3, _ := future3.Await()

// All three operations run concurrently, completing in ~100ms total
```

### Error Handling

```go
future := async.Async(ctx, func(ctx context.Context) (int, error) {
    return 0, errors.New("something went wrong")
})

result, err := future.Await()
if err != nil {
    fmt.Println("Error:", err)
}
```

### Panic Recovery

The package automatically recovers from panics and converts them to errors:

```go
future := async.Async(ctx, func(ctx context.Context) (int, error) {
    panic("unexpected error")
})

result, err := future.Await()
// err will contain the panic message with stack trace
```

### Multiple Awaits

You can call `Await()` multiple times on the same Future. The function is only executed once, and subsequent calls return the cached result:

```go
future := async.Async(ctx, func(ctx context.Context) (int, error) {
    return 42, nil
})

result1, _ := future.Await() // Executes the function
result2, _ := future.Await() // Returns cached result
```

## API Reference

### Types

#### `Future[T any]`

Represents a value that will be available at some point in the future.

#### `ErrCancelled`

A predefined error that can be used to signal cancellation.

### Functions

#### `Async[T any](ctx context.Context, f func(context.Context) (T, error)) *Future[T]`

Wraps a function and executes it asynchronously in a goroutine. Returns a `Future[T]` that can be used to retrieve the result.

**Parameters:**
- `ctx`: Context for cancellation and value propagation
- `f`: Function to execute asynchronously

**Returns:**
- A pointer to a `Future[T]` containing the eventual result

#### `(*Future[T]) Await() (T, error)`

Waits for the Future to complete and returns the result. Blocks until the asynchronous operation finishes.

**Returns:**
- The result value of type `T`
- Any error that occurred during execution (including recovered panics)

## License

MIT
