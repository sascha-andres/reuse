package async

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAsyncAwaitBasic(t *testing.T) {
	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (int, error) {
		return 42, nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != 42 {
		t.Fatalf("expected 42, got: %d", result)
	}
}

func TestAsyncAwaitString(t *testing.T) {
	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (string, error) {
		return "hello world", nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != "hello world" {
		t.Fatalf("expected 'hello world', got: %s", result)
	}
}

func TestAsyncAwaitStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (Person, error) {
		return Person{Name: "Alice", Age: 30}, nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Name != "Alice" || result.Age != 30 {
		t.Fatalf("expected Person{Name: 'Alice', Age: 30}, got: %+v", result)
	}
}

func TestAsyncAwaitError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("test error")

	future := Async(ctx, func(ctx context.Context) (int, error) {
		return 0, expectedErr
	})

	result, err := future.Await()
	if err != expectedErr {
		t.Fatalf("expected error %v, got: %v", expectedErr, err)
	}
	if result != 0 {
		t.Fatalf("expected zero value, got: %d", result)
	}
}

func TestAsyncAwaitPanicWithError(t *testing.T) {
	ctx := context.Background()
	panicErr := errors.New("panic error")

	future := Async(ctx, func(ctx context.Context) (int, error) {
		panic(panicErr)
	})

	result, err := future.Await()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "recovering from error") {
		t.Fatalf("expected error to contain 'recovering from error', got: %v", err)
	}
	if !strings.Contains(err.Error(), "panic error") {
		t.Fatalf("expected error to contain 'panic error', got: %v", err)
	}
	if result != 0 {
		t.Fatalf("expected zero value, got: %d", result)
	}
}

func TestAsyncAwaitPanicWithString(t *testing.T) {
	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (int, error) {
		panic("string panic")
	})

	result, err := future.Await()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "panic:") {
		t.Fatalf("expected error to contain 'panic:', got: %v", err)
	}
	if !strings.Contains(err.Error(), "string panic") {
		t.Fatalf("expected error to contain 'string panic', got: %v", err)
	}
	if result != 0 {
		t.Fatalf("expected zero value, got: %d", result)
	}
}

func TestAsyncAwaitMultipleAwaits(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	future := Async(ctx, func(ctx context.Context) (int, error) {
		callCount++
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	})

	result1, err1 := future.Await()
	if err1 != nil {
		t.Fatalf("expected no error on first await, got: %v", err1)
	}
	if result1 != 42 {
		t.Fatalf("expected 42 on first await, got: %d", result1)
	}

	result2, err2 := future.Await()
	if err2 != nil {
		t.Fatalf("expected no error on second await, got: %v", err2)
	}
	if result2 != 42 {
		t.Fatalf("expected 42 on second await, got: %d", result2)
	}

	if callCount != 1 {
		t.Fatalf("expected function to be called once, was called %d times", callCount)
	}
}

func TestAsyncAwaitWithDelay(t *testing.T) {
	ctx := context.Background()

	startTime := time.Now()
	future := Async(ctx, func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "delayed result", nil
	})

	result, err := future.Await()
	elapsed := time.Since(startTime)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != "delayed result" {
		t.Fatalf("expected 'delayed result', got: %s", result)
	}
	if elapsed < 200*time.Millisecond {
		t.Fatalf("expected at least 200ms delay, got: %v", elapsed)
	}
}

func TestAsyncAwaitConcurrent(t *testing.T) {
	ctx := context.Background()

	future1 := Async(ctx, func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 1, nil
	})

	future2 := Async(ctx, func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 2, nil
	})

	future3 := Async(ctx, func(ctx context.Context) (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 3, nil
	})

	startTime := time.Now()

	result1, err1 := future1.Await()
	result2, err2 := future2.Await()
	result3, err3 := future3.Await()

	elapsed := time.Since(startTime)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("expected no errors, got: %v, %v, %v", err1, err2, err3)
	}

	if result1 != 1 || result2 != 2 || result3 != 3 {
		t.Fatalf("expected results 1, 2, 3, got: %d, %d, %d", result1, result2, result3)
	}

	if elapsed >= 300*time.Millisecond {
		t.Fatalf("expected concurrent execution (< 300ms), got: %v", elapsed)
	}
}

func TestAsyncAwaitPointer(t *testing.T) {
	type Data struct {
		Value int
	}

	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (*Data, error) {
		return &Data{Value: 100}, nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatalf("expected non-nil result")
	}
	if result.Value != 100 {
		t.Fatalf("expected Value 100, got: %d", result.Value)
	}
}

func TestAsyncAwaitSlice(t *testing.T) {
	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) ([]int, error) {
		return []int{1, 2, 3, 4, 5}, nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result) != 5 {
		t.Fatalf("expected length 5, got: %d", len(result))
	}
	for i, v := range result {
		if v != i+1 {
			t.Fatalf("expected value %d at index %d, got: %d", i+1, i, v)
		}
	}
}

func TestAsyncAwaitMap(t *testing.T) {
	ctx := context.Background()

	future := Async(ctx, func(ctx context.Context) (map[string]int, error) {
		return map[string]int{"one": 1, "two": 2}, nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected length 2, got: %d", len(result))
	}
	if result["one"] != 1 || result["two"] != 2 {
		t.Fatalf("unexpected map values: %v", result)
	}
}

func TestAsyncWithContextValues(t *testing.T) {
	type contextKey string
	key := contextKey("testKey")

	ctx := context.WithValue(context.Background(), key, "testValue")

	future := Async(ctx, func(ctx context.Context) (string, error) {
		value := ctx.Value(key)
		if value == nil {
			return "", fmt.Errorf("context value not found")
		}
		return value.(string), nil
	})

	result, err := future.Await()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != "testValue" {
		t.Fatalf("expected 'testValue', got: %s", result)
	}
}

func TestErrCancelled(t *testing.T) {
	if ErrCancelled == nil {
		t.Fatalf("ErrCancelled should not be nil")
	}
	if ErrCancelled.Error() != "cancelled" {
		t.Fatalf("expected error message 'cancelled', got: %s", ErrCancelled.Error())
	}
}
