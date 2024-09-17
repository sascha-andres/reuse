package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// CircuitBreaker states
const (

	// StateClosed represents the closed state of a CircuitBreaker, where requests are allowed to pass through.
	StateClosed = "Closed"

	// StateOpen represents the open state of a CircuitBreaker, where requests are blocked.
	StateOpen = "Open"

	// StateHalfOpen represents the half-open state of a CircuitBreaker, where limited requests are allowed to test recovery.
	StateHalfOpen = "Half-Open"
)

// CircuitBreaker struct
type CircuitBreaker struct {

	// mu is a mutual exclusion lock to protect critical sections within the CircuitBreaker.
	mu sync.Mutex

	// failures track the number of consecutive task failures to manage the circuit breaker states.
	failures int

	// state represents the current state of the circuit breaker (e.g., Open, Half-Open, Closed).
	state string

	// maxFailures is the maximum number of consecutive task failures allowed before the circuit breaker transitions to the open state.
	maxFailures int

	// timeout is the duration the circuit breaker waits before transitioning from an open to a half-open state.
	timeout time.Duration

	// lastFailureTime denotes the time at which the last task failure occurred, used to manage circuit breaker state transitions.
	lastFailureTime time.Time

	// fn is the function executed by the CircuitBreaker, returning an error if the task fails.
	fn WorkFunc
}

// NewCircuitBreaker initializes and returns a new CircuitBreaker instance with specified maxFailures, timeout, and function.
func NewCircuitBreaker(maxFailures int, timeout time.Duration, fn WorkFunc) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		maxFailures: maxFailures,
		timeout:     timeout,
		fn:          fn,
	}
}

// Call executes a task within the CircuitBreaker, transitioning states based on task success or failure.
func (cb *CircuitBreaker) Call(wg *sync.WaitGroup, taskDone chan<- Task, id int) {
	defer wg.Done()
	cb.mu.Lock()

	//  circuit breaker state
	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
		} else {
			cb.mu.Unlock()
			taskDone <- Task{Id: id}
			return
		}
	case StateHalfOpen:
		// Allow some tasks to test if the service has recovered
		// continue with the instructions
	}

	cb.mu.Unlock()

	// Perform the task
	err := cb.fn(id)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		fmt.Println("Task failed. Failure count:", cb.failures)

		// If max failures reached, open the circuit
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			cb.lastFailureTime = time.Now()
		}
		taskDone <- Task{Id: id}
	} else {
		// Success: reset failure count
		cb.failures = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
		}
		taskDone <- Task{Id: id, Status: true}
	}
}

// WorkFunc defines a function type that takes an integer ID and returns an error.
type WorkFunc func(id int) error

// Task represents a unit of work with an ID and completion status.
type Task struct {

	// Status indicates whether the task was successfully completed.
	Status bool

	// Id is the unique identifier for a task.
	Id int
}
