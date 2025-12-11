package utils

import (
	"context"
	"fmt"
	"time"
)

// Promise represents a future result of type T with context support
type Promise[T any] struct {
	result T
	err    error
	done   chan struct{}
	cancel context.CancelFunc
	ctx    context.Context
}

// Async runs a function in a separate goroutine and returns a Promise.
// The provided function receives a context that is cancelled when Promise.Cancel() is called.
func Async[T any](fn func(ctx context.Context) (T, error)) *Promise[T] {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Promise[T]{
		done:   make(chan struct{}),
		cancel: cancel,
		ctx:    ctx,
	}

	go func() {
		defer close(p.done)
		defer func() {
			if r := recover(); r != nil {
				p.err = fmt.Errorf("panic in Async task: %v", r)
			}
		}()

		p.result, p.err = fn(ctx)
	}()

	return p
}

// AsyncWithContext runs a function with a parent context
func AsyncWithContext[T any](parentCtx context.Context, fn func(ctx context.Context) (T, error)) *Promise[T] {
	ctx, cancel := context.WithCancel(parentCtx)
	p := &Promise[T]{
		done:   make(chan struct{}),
		cancel: cancel,
		ctx:    ctx,
	}

	go func() {
		defer close(p.done)
		defer func() {
			if r := recover(); r != nil {
				p.err = fmt.Errorf("panic in Async task: %v", r)
			}
		}()

		p.result, p.err = fn(ctx)
	}()

	return p
}

// Await blocks until the promise is resolved and returns the result
func (p *Promise[T]) Await() (T, error) {
	<-p.done
	return p.result, p.err
}

// Cancel cancels the underlying context of the async task
func (p *Promise[T]) Cancel() {
	p.cancel()
}

// WithTimeout sets a timeout for the promise. If it expires, the context is cancelled.
// Note: This does NOT wait for the task to finish, it just signals cancellation.
func (p *Promise[T]) WithTimeout(d time.Duration) *Promise[T] {
	// We create a new context with timeout derived from the promise's current context
	// However, since the goroutine is already running with p.ctx, we can't easily swap it.
	// Instead, we just schedule a cancel after duration.
	time.AfterFunc(d, func() {
		p.cancel()
	})
	return p
}

// PromiseAll waits for all promises to complete
func PromiseAll[T any](promises ...*Promise[T]) ([]T, error) {
	results := make([]T, len(promises))
	for i, p := range promises {
		res, err := p.Await()
		if err != nil {
			return nil, err
		}
		results[i] = res
	}
	return results, nil
}

// PromiseRace returns the result of the first promise to complete
func PromiseRace[T any](promises ...*Promise[T]) (T, error) {
	if len(promises) == 0 {
		var zero T
		return zero, fmt.Errorf("no promises to race")
	}

	resChan := make(chan T)
	errChan := make(chan error)
	done := make(chan struct{})

	// Ensure we only close channels once
	defer close(done)

	for _, p := range promises {
		go func(promise *Promise[T]) {
			// Wait for this promise
			res, err := promise.Await()

			// Try to send result, but don't block if race is over
			select {
			case <-done:
				return
			default:
				if err != nil {
					select {
					case errChan <- err:
					case <-done:
					}
				} else {
					select {
					case resChan <- res:
					case <-done:
					}
				}
			}
		}(p)
	}

	select {
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		var zero T
		return zero, err
	}
}

// PromiseResult holds the result of a PromiseAllSettled call
type PromiseResult[T any] struct {
	Status string // "fulfilled" or "rejected"
	Value  T
	Error  error
}

// PromiseAllSettled waits for all promises to finish, regardless of success/failure
func PromiseAllSettled[T any](promises ...*Promise[T]) []PromiseResult[T] {
	results := make([]PromiseResult[T], len(promises))
	for i, p := range promises {
		res, err := p.Await()
		if err != nil {
			results[i] = PromiseResult[T]{Status: "rejected", Error: err}
		} else {
			results[i] = PromiseResult[T]{Status: "fulfilled", Value: res}
		}
	}
	return results
}

// GlobalErrorHandler is called when a background job fails
var GlobalErrorHandler = func(err error) {
	fmt.Printf("Background job error: %v\n", err)
}

// RunBackground runs a task in the background (fire and forget)
func RunBackground(fn func() error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				GlobalErrorHandler(fmt.Errorf("panic in background job: %v", r))
			}
		}()
		if err := fn(); err != nil {
			GlobalErrorHandler(err)
		}
	}()
}
