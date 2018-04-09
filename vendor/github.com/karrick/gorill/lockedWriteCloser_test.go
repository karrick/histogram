package gorill

import (
	"io"
	"testing"
)

func BenchmarkWriterLockedWriter(b *testing.B) {
	consumers := make([]io.WriteCloser, consumerCount)
	for i := 0; i < len(consumers); i++ {
		consumers[i] = NewLockingWriteCloser(NewNopCloseBuffer())
	}
	benchmarkWriter(b, b.N, consumers)
}
