package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// CircuitBreaker states
const (
	StateClosed   = "Closed"
	StateOpen     = "Open"
	StateHalfOpen = "Half-Open"
)

// CircuitBreaker struct
type CircuitBreaker struct {
	mu              sync.Mutex
	failures        int
	state           string
	maxFailures     int
	timeout         time.Duration
	lastFailureTime time.Time
}

// initializes the circuit breaker
func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:       StateClosed,
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

// simulates executing tasks and handling circuit breaker states
func (cb *CircuitBreaker) Call(wg *sync.WaitGroup, taskDone chan<- Task, id int) {
	defer wg.Done()
	cb.mu.Lock()

	//  circuit breaker state
	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			fmt.Println("Circuit transitioning to Half-Open...")
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
	err := SimulateWork(id)

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		fmt.Println("Task failed. Failure count:", cb.failures)

		// If max failures reached, open the circuit
		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
			cb.lastFailureTime = time.Now()
			fmt.Println("Developer lost tempo... needs a full break!")
		}
		taskDone <- Task{Id: id}
	} else {
		// Success: reset failure count
		cb.failures = 0
		if cb.state == StateHalfOpen {
			cb.state = StateClosed
			fmt.Println("Developer is back from break and open to accept new tasks ")
		}
		taskDone <- Task{Id: id, Status: true}
	}
}
