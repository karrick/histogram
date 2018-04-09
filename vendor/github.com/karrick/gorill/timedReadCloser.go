package gorill

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// TimedReadCloser is an io.Reader that enforces a preset timeout period on every Read operation.
type TimedReadCloser struct {
	halted   bool
	iorc     io.ReadCloser
	jobs     chan *rillJob
	jobsDone sync.WaitGroup
	lock     sync.RWMutex
	timeout  time.Duration
}

// NewTimedReadCloser returns a TimedReadCloser that enforces a preset timeout period on every Read
// operation.  It panics when timeout is less than or equal to 0.
func NewTimedReadCloser(iowc io.ReadCloser, timeout time.Duration) *TimedReadCloser {
	if timeout <= 0 {
		panic(fmt.Errorf("timeout must be greater than 0: %s", timeout))
	}
	rc := &TimedReadCloser{
		iorc:    iowc,
		jobs:    make(chan *rillJob, 1),
		timeout: timeout,
	}
	rc.jobsDone.Add(1)
	go func() {
		for job := range rc.jobs {
			n, err := rc.iorc.Read(job.data)
			job.results <- rillResult{n, err}
		}
		rc.jobsDone.Done()
	}()
	return rc
}

// Read reads data to the underlying io.Reader, but returns ErrTimeout if the Read operation exceeds
// a preset timeout duration.
//
// Even after a timeout takes place, the read may still independently complete as reads are queued
// from a different go-routine.  Race condition for the data slice is prevented by reading into a
// temporary byte slice, and copying the results to the client's slice when the actual read returns.
func (rc *TimedReadCloser) Read(data []byte) (int, error) {
	rc.lock.RLock()
	defer rc.lock.RUnlock()

	if rc.halted {
		return 0, ErrReadAfterClose{}
	}

	job := newRillJob(_read, make([]byte, len(data)))
	rc.jobs <- job

	// wait for result or timeout
	select {
	case result := <-job.results:
		copy(data, job.data)
		return result.n, result.err
	case <-time.After(rc.timeout):
		return 0, ErrTimeout(rc.timeout)
	}
}

// Close frees resources when a SpooledReadCloser is no longer needed.
func (rc *TimedReadCloser) Close() error {
	rc.lock.Lock()
	defer rc.lock.Unlock()

	close(rc.jobs)
	rc.jobsDone.Wait()
	rc.halted = true
	return rc.iorc.Close()
}
