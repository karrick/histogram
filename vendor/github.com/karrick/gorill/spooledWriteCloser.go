package gorill

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"
)

// DefaultBufSize is the default size of the underlying bufio.Writer buffer.
const DefaultBufSize = 4096

// DefaultFlushPeriod is the default frequency of buffer flushes.
const DefaultFlushPeriod = 15 * time.Second

// SpooledWriteCloser spools bytes written to it through a bufio.Writer, periodically flushing data
// written to underlying io.WriteCloser.
type SpooledWriteCloser struct {
	bufSize     int
	bw          *bufio.Writer
	flushPeriod time.Duration
	halted      bool
	iowc        io.WriteCloser
	jobs        chan *rillJob
	jobsDone    sync.WaitGroup
	lock        sync.RWMutex
}

// SpooledWriteCloserSetter is any function that modifies a SpooledWriteCloser being instantiated.
type SpooledWriteCloserSetter func(*SpooledWriteCloser) error

// Flush is used to configure a new SpooledWriteCloser to periodically flush.
func Flush(periodicity time.Duration) SpooledWriteCloserSetter {
	return func(sw *SpooledWriteCloser) error {
		if periodicity <= 0 {
			return fmt.Errorf("periodicity must be greater than 0: %s", periodicity)
		}
		sw.flushPeriod = periodicity
		return nil
	}
}

// BufSize is used to configure a new SpooledWriteCloser's buffer size.
func BufSize(size int) SpooledWriteCloserSetter {
	return func(sw *SpooledWriteCloser) error {
		if size <= 0 {
			return fmt.Errorf("buffer size must be greater than 0: %d", size)
		}
		sw.bufSize = size
		return nil
	}
}

// NewSpooledWriteCloser returns a SpooledWriteCloser that spools bytes written to it through a
// bufio.Writer, periodically forcing the bufio.Writer to flush its contents.
func NewSpooledWriteCloser(iowc io.WriteCloser, setters ...SpooledWriteCloserSetter) (*SpooledWriteCloser, error) {
	w := &SpooledWriteCloser{
		bufSize:     DefaultBufSize,
		flushPeriod: DefaultFlushPeriod,
		iowc:        iowc,
		jobs:        make(chan *rillJob, 1),
	}
	for _, setter := range setters {
		if err := setter(w); err != nil {
			return nil, err
		}
	}
	w.bw = bufio.NewWriterSize(iowc, w.bufSize)
	w.jobsDone.Add(1)
	go func() {
		ticker := time.NewTicker(w.flushPeriod)
		defer ticker.Stop()
		defer w.jobsDone.Done()
		for {
			select {
			case job, more := <-w.jobs:
				if !more {
					return
				}
				switch job.op {
				case _write:
					n, err := w.bw.Write(job.data)
					job.results <- rillResult{n, err}
				case _flush:
					err := w.bw.Flush()
					job.results <- rillResult{0, err}
				}
			case <-ticker.C:
				w.bw.Flush()
			}
		}
	}()
	return w, nil
}

// Write spools a byte slice of data to be written to the SpooledWriteCloser.
func (w *SpooledWriteCloser) Write(data []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.halted {
		return 0, ErrWriteAfterClose{}
	}

	job := newRillJob(_write, data)
	w.jobs <- job
	// wait for results
	result := <-job.results
	return result.n, result.err
}

// Flush causes all data not yet written to the output stream to be flushed.
func (w *SpooledWriteCloser) Flush() error {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.halted {
		return ErrWriteAfterClose{}
	}

	job := newRillJob(_flush, nil)
	w.jobs <- job
	result := <-job.results
	// wait for results
	return result.err
}

// Close frees resources when a SpooledWriteCloser is no longer needed.
func (w *SpooledWriteCloser) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	close(w.jobs)
	w.jobsDone.Wait()
	w.halted = true

	var errors ErrList
	errors.Append(w.bw.Flush())
	errors.Append(w.iowc.Close())
	return errors.Err()
}
