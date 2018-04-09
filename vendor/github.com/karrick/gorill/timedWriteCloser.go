package gorill

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// ErrTimeout error is returned whenever a Write operation exceeds the preset timeout period. Even
// after a timeout takes place, the write may still independantly complete.
type ErrTimeout time.Duration

// Error returns a string representing the ErrTimeout.
func (e ErrTimeout) Error() string {
	return fmt.Sprintf("timeout after %s", time.Duration(e))
}

// TimedWriteCloser is an io.Writer that enforces a preset timeout period on every Write operation.
type TimedWriteCloser struct {
	halted   bool
	iowc     io.WriteCloser
	jobs     chan *rillJob
	jobsDone sync.WaitGroup
	lock     sync.RWMutex
	timeout  time.Duration
}

// NewTimedWriteCloser returns a TimedWriteCloser that enforces a preset timeout period on every Write
// operation.  It panics when timeout is less than or equal to 0.
func NewTimedWriteCloser(iowc io.WriteCloser, timeout time.Duration) *TimedWriteCloser {
	if timeout <= 0 {
		panic(fmt.Errorf("timeout must be greater than 0: %s", timeout))
	}
	wc := &TimedWriteCloser{
		iowc:    iowc,
		jobs:    make(chan *rillJob, 1),
		timeout: timeout,
	}
	wc.jobsDone.Add(1)
	go func() {
		for job := range wc.jobs {
			n, err := wc.iowc.Write(job.data)
			job.results <- rillResult{n, err}
		}
		wc.jobsDone.Done()
	}()
	return wc
}

// Write writes data to the underlying io.Writer, but returns ErrTimeout if the Write
// operation exceeds a preset timeout duration.  Even after a timeout takes place, the write may
// still independantly complete as writes are queued from a different go routine.
func (wc *TimedWriteCloser) Write(data []byte) (int, error) {
	wc.lock.RLock()
	defer wc.lock.RUnlock()

	if wc.halted {
		return 0, ErrWriteAfterClose{}
	}

	job := newRillJob(_write, data)
	wc.jobs <- job

	// wait for result or timeout
	select {
	case result := <-job.results:
		return result.n, result.err
	case <-time.After(wc.timeout):
		return 0, ErrTimeout(wc.timeout)
	}
}

// Close frees resources when a SpooledWriteCloser is no longer needed.
func (wc *TimedWriteCloser) Close() error {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	close(wc.jobs)
	wc.jobsDone.Wait()
	wc.halted = true
	return wc.iowc.Close()
}
