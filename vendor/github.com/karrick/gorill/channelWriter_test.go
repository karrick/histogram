package gorill

import (
	"io"
	"sync"
	"testing"
)

// channelWriter provided to benchmark against LockingWriter and TimedWriteCloser.
type channelWriter struct {
	halted   bool
	jobsDone sync.WaitGroup
	iowc     io.WriteCloser
	jobs     chan rillJob
	lock     sync.RWMutex
}

func newChannelWriter(iowc io.WriteCloser) *channelWriter {
	w := &channelWriter{
		iowc: iowc,
		jobs: make(chan rillJob, 1),
	}
	go func(w *channelWriter) {
		w.jobsDone.Add(1)
		for job := range w.jobs {
			n, err := w.iowc.Write(job.data)
			job.results <- rillResult{n, err}
		}
		w.jobsDone.Done()
	}(w)
	return w
}

func (w *channelWriter) Write(data []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.halted {
		return 0, ErrWriteAfterClose{}
	}

	job := rillJob{data: data, results: make(chan rillResult, 1)}
	w.jobs <- job
	// wait for result
	result := <-job.results
	return result.n, result.err
}

func (w *channelWriter) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	close(w.jobs)
	w.jobsDone.Wait()
	w.halted = true
	return w.iowc.Close()
}

func BenchmarkWriterChannelWriter(b *testing.B) {
	consumers := make([]io.WriteCloser, consumerCount)
	for i := 0; i < len(consumers); i++ {
		consumers[i] = newChannelWriter(NewNopCloseBuffer())
	}
	benchmarkWriter(b, b.N, consumers)
}
