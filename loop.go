package concurrent

import (
	"fmt"
	"runtime"
	"sync"
)

// An ExecutionError indicates that some calling went wrong.
type ExecutionError struct {
	Cause      error   // First encountered error
	Suppressed []error // Subsequent but suppressed errors
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("cannot execute: %v", e.Cause)
}

func (e ExecutionError) Unwrap() error {
	return e.Cause
}

// ForEach is a simple alias to Execute and just spawns runtime.NumCPU() routines.
func ForEach(length int, f func(idx int) error) error {
	return Execute(runtime.NumCPU(), length, nil, f)
}

// Execute blocks and executes with at most n-routines about length indices and calls f with
// the index in no defined order.
// The first encountered error is returned and pending routines are tried to get cancelled.
//
// The execution is tried to be aborted, if cancel is set. Why no channel here? Because a channel is actually
// more complicated:
//  * it is slower, because it uses a full memory barrier (mutex), see src/runtime/chan.go
//  * it is hard to use right and reason about: lifetime and resource leaks? how many cancel elements to consume? Who
//    closes whom?
// If any error occurs, the execution is tried to be cancelled early and an ExecutionError is returned.
func Execute(n, length int, cancel *AtomicBool, f func(idx int) error) error {
	// if no length, nothing to do or to allocate at all
	if length <= 0 {
		return nil
	}

	// if just one element, execute immediately without overhead
	if length == 1 {
		return f(0)
	}

	// if more than length, no need to spawn more routines
	if n > length {
		n = length
	}

	// fix bogus concurrency (e.g. 1 core / 2 = 0)
	if n < 1 {
		n = 1
	}

	queue := make(chan int, length)
	for i := 0; i < length; i++ {
		queue <- i
	}
	var errs Vec

	defer close(queue)

	wg := &sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case next := <-queue:
					if cancel.Load() || errs.Len() > 0 {
						return
					}

					err := f(next)
					if err != nil {
						errs.PushBack(err)
						return
					}

				default:
					return
				}
			}

		}(i)
	}

	wg.Wait()

	if firstErr, ok := errs.PopFirst(); ok {
		var err ExecutionError
		err.Cause = firstErr.(error)
		for otherErr, ok := errs.PopBack(); ok; {
			err.Suppressed = append(err.Suppressed, otherErr.(error))
		}

		return err
	}

	return nil
}
