// +build !race

package gorill

import (
	"io"
	"testing"
	"time"
)

func TestMultiWriteCloserFanIn(t *testing.T) {
	bb := NewNopCloseBufferSize(16384)

	first := NewMultiWriteCloserFanIn(bb)

	want := string(largeBuf)
	first.Write(largeBuf)

	if actual := bb.String(); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}

	bb.Reset()
	want = ""
	if actual := bb.String(); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}

	second := first.Add()
	want = string(largeBuf)
	second.Write(largeBuf)

	if actual := bb.String(); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}

	first.Close()
	if want, actual := false, bb.IsClosed(); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}

	second.Close()
	time.Sleep(100 * time.Millisecond) // race condition during testing
	if want, actual := true, bb.IsClosed(); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
}

func BenchmarkWriterMultiWriteCloserFanIn(b *testing.B) {
	consumers := make([]io.WriteCloser, consumerCount)
	for i := 0; i < len(consumers); i++ {
		consumers[i] = NewMultiWriteCloserFanIn(NewNopCloseBuffer())
	}
	benchmarkWriter(b, b.N, consumers)
}
