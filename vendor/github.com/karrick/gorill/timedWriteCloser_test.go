package gorill

import (
	"bytes"
	"io"
	"testing"
	"time"
)

var timedWriterBuf []byte

func init() {
	timedWriterBuf = make([]byte, 1024)
	for i := range timedWriterBuf {
		timedWriterBuf[i] = '.'
	}
}

func TestTimedWriteCloserBeforeTimeout(t *testing.T) {
	timeout := time.Second
	bb := new(bytes.Buffer)

	tw := NewTimedWriteCloser(NopCloseWriter(SlowWriter(bb, timeout/10)), timeout)
	defer tw.Close()

	n, err := tw.Write(timedWriterBuf)
	if want := len(timedWriterBuf); n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := string(timedWriterBuf); want != bb.String() {
		t.Errorf("Actual: %#v; Expected: %#v", bb.String(), want)
	}
}

func TestTimedWriteCloserAfterTimeout(t *testing.T) {
	timeout := time.Millisecond
	bb := new(bytes.Buffer)

	tw := NewTimedWriteCloser(NopCloseWriter(SlowWriter(bb, 10*timeout)), timeout)
	defer tw.Close()

	n, err := tw.Write(timedWriterBuf)
	if want := 0; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if _, ok := err.(ErrTimeout); err == nil || !ok {
		t.Errorf("Actual: %#v; Expected: %s", err, ErrTimeout(timeout))
	}
	// NOTE: cannot check for contents of buffer, because write independently completes.
}

func TestTimedWriteCloserWriteAfterCloseReturnsError(t *testing.T) {
	tw := NewTimedWriteCloser(NopCloseWriter(new(bytes.Buffer)), time.Millisecond)

	tw.Close()

	n, err := tw.Write([]byte(alphabet))
	if want := 0; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if _, ok := err.(ErrWriteAfterClose); err == nil || !ok {
		t.Errorf("Actual: %s; Expected: %#v", err, ErrWriteAfterClose{})
	}
}

func BenchmarkWriterTimedWriteCloser(b *testing.B) {
	consumers := make([]io.WriteCloser, consumerCount)
	for i := 0; i < len(consumers); i++ {
		consumers[i] = NewTimedWriteCloser(NewNopCloseBuffer(), time.Minute)
	}
	benchmarkWriter(b, b.N, consumers)
}
