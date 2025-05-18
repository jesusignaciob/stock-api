package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// ZeroValue returns the zero value for any type T.
// This is useful for returning a default value in case of errors or timeouts.
//
// Type Parameters:
//   - T: Any type.
//
// Returns:
//   - The zero value of type T.
func ZeroValue[T any]() T {
	var zero T
	return zero
}

// AsyncResult is a generic struct that holds the result of an asynchronous operation.
// Fields:
//   - Result: The result of type T returned by the operation.
//   - Count: An integer count, useful for operations that return a count (e.g., number of records).
//   - Error: Any error encountered during the operation.
type AsyncResult[T any] struct {
	Result T
	Count  int
	Error  error
}

// AsyncOperation executes the provided operation asynchronously using a worker pool.
// It returns the result of the operation or an error if the operation times out, the client disconnects,
// or the server is busy (i.e., the worker pool is full).
//
// Type Parameters:
//   - T: The type of the result returned by the operation.
//
// Parameters:
//   - c: The Gin context, used to detect client disconnects.
//   - workerPool: A channel used to limit the number of concurrent operations.
//   - operation: A function that performs the operation and returns a result of type T and an error.
//
// Returns:
//   - T: The result of the operation, or the zero value of T in case of error.
//   - error: An error if the operation fails, times out, the client disconnects, or the server is busy.
//
// The function waits up to 5 seconds for the operation to complete. If the operation does not complete
// within this time, it returns a timeout error. If the client disconnects before the operation completes,
// it returns a client disconnected error. If the worker pool is full, it returns a server busy error.
func AsyncOperation[T any](
	c *gin.Context,
	workerPool chan struct{},
	operation func() (T, error),
) (T, error) {
	resultChan := make(chan AsyncResult[T], 1)

	select {
	case workerPool <- struct{}{}:
		go func() {
			defer func() { <-workerPool }()

			result, err := operation()
			resultChan <- AsyncResult[T]{Result: result, Error: err}
		}()

		select {
		case res := <-resultChan:
			return res.Result, res.Error
		case <-time.After(5 * time.Second):
			return ZeroValue[T](), fmt.Errorf("operation timeout")
		case <-c.Request.Context().Done():
			return ZeroValue[T](), fmt.Errorf("client disconnected")
		}
	default:
		return ZeroValue[T](), fmt.Errorf("server busy")
	}
}

// AsyncManyOperation executes the provided operation asynchronously using a worker pool,
// and returns its result, count, and error. It leverages Go generics to support any result type.
// The function ensures that the number of concurrent operations does not exceed the worker pool capacity.
// It waits for the operation to complete, a timeout (5 seconds), or client disconnection, whichever comes first.
//
// Parameters:
//   - c: The Gin context, used to detect client disconnection.
//   - workerPool: A channel used to limit the number of concurrent operations.
//   - operation: A function that performs the desired operation and returns a result of type T, a count, and an error.
//
// Returns:
//   - T: The result of the operation (zero value if failed or timed out).
//   - int: The count returned by the operation (0 if failed or timed out).
//   - error: An error if the operation failed, timed out, the client disconnected, or the server is busy.
//
// Possible errors:
//   - "operation timeout": If the operation does not complete within 5 seconds.
//   - "client disconnected": If the client disconnects before the operation completes.
//   - "server busy": If the worker pool is full and cannot accept new operations.
func AsyncManyOperation[T any](
	c *gin.Context,
	workerPool chan struct{},
	operation func() (T, int, error),
) (result T, count int, err error) {
	resultChan := make(chan AsyncResult[T], 1)

	select {
	case workerPool <- struct{}{}:
		go func() {
			defer func() { <-workerPool }()

			result, count, err := operation()
			resultChan <- AsyncResult[T]{Result: result, Count: count, Error: err}
		}()

		select {
		case res := <-resultChan:
			return res.Result, res.Count, res.Error
		case <-time.After(5 * time.Second):
			return ZeroValue[T](), 0, fmt.Errorf("operation timeout")
		case <-c.Request.Context().Done():
			return ZeroValue[T](), 0, fmt.Errorf("client disconnected")
		}
	default:
		return ZeroValue[T](), 0, fmt.Errorf("server busy")
	}
}
