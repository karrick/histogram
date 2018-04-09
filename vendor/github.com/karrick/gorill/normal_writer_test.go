// +build !race

package gorill

import (
	"io"
	"testing"
)

// The test contains a data race to determine maximum throughput of writing, while ignoring race
// conditions.
func BenchmarkWriterNormalWriter(b *testing.B) {
	consumers := make([]io.WriteCloser, consumerCount)
	for i := 0; i < len(consumers); i++ {
		consumers[i] = NewNopCloseBuffer()
	}
	benchmarkWriter(b, b.N, consumers)
}
