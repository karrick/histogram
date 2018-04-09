package gorill

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"testing"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz\n"
const consumerCount = 1000

func benchmarkWriter(b *testing.B, n int, consumers []io.WriteCloser) {
	return // FIXME

	const producerCount = 10

	producerBegin := make(chan struct{})
	var finished sync.WaitGroup

	finished.Add(producerCount)

	for i := 0; i < producerCount; i++ {
		go func(begin <-chan struct{}, finished *sync.WaitGroup) {
			<-begin
			for i := 0; i < n; i++ {
				j := rand.Int31n(int32(consumerCount))
				n, err := consumers[j].Write([]byte(alphabet))
				if want := len(alphabet); n != want {
					b.Fatalf("Actual: %#v; Expected: %#v", n, want)
				}
				if err != nil {
					b.Fatalf("Actual: %#v; Expected: %#v", err, nil)
				}
			}
			finished.Done()
		}(producerBegin, &finished)
	}

	b.ResetTimer()
	for i := 0; i < producerCount; i++ {
		producerBegin <- struct{}{}
	}
	fmt.Printf("\nwaiting for producer complete: ")
	finished.Wait()
	fmt.Printf("\ndone waiting for producer complete")
	for i := 0; i < consumerCount; i++ {
		consumers[i].Close()
	}
}
