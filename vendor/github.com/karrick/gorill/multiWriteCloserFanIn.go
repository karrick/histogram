package gorill

import (
	"io"
	"sync"
)

const (
	bufSize = 4096
)

// MultiWriteCloserFanIn is a structure that provides multiple io.WriteClosers to write to same underlying
// io.WriteCloser.  When the final io.WriteCloser that MultiWriteCloserFanIn provides is closed, then the
// underlying io.WriteCloser will be closed.
type MultiWriteCloserFanIn struct {
	iowc  io.WriteCloser
	done  sync.WaitGroup
	pLock *sync.Mutex
	pDone *sync.WaitGroup
}

// NewMultiWriteCloserFanIn creates a MultiWriteCloserFanIn instance where writes to any of the provided
// io.WriteCloser instances will be funneled to the underlying io.WriteCloser instance.  The client
// ought to call Close on all provided io.WriteCloser instances, after which, MultiWriteCloserFanIn will
// close the underlying io.WriteCloser.
//
//    func Example(largeBuf []byte) {
//    	bb := NewNopCloseBufferSize(16384)
//    	first := NewMultiWriteCloserFanIn(bb)
//    	second := first.Add()
//    	first.Write(largeBuf)
//    	first.Close()
//    	second.Write(largeBuf)
//    	second.Close()
//    }
func NewMultiWriteCloserFanIn(iowc io.WriteCloser) *MultiWriteCloserFanIn {
	var lock sync.Mutex
	var done sync.WaitGroup
	prime := &MultiWriteCloserFanIn{iowc: iowc, pLock: &lock, pDone: &done}
	d := prime.Add()
	go func() {
		done.Wait()
		iowc.Close()
	}()
	return d
}

// Add returns a new MultiWriteCloserFanIn that redirects all writes to the underlying
// io.WriteCloser.  The client ought to call Close on the returned MultiWriteCloserFanIn to signify
// intent to no longer Write to the MultiWriteCloserFanIn.
func (fanin *MultiWriteCloserFanIn) Add() *MultiWriteCloserFanIn {
	d := &MultiWriteCloserFanIn{iowc: fanin.iowc, pLock: fanin.pLock, pDone: fanin.pDone}
	d.done.Add(1)
	d.pDone.Add(1)
	go func() {
		d.done.Wait()
		d.pDone.Done()
	}()
	return d
}

// Write copies the entire data slice to the underlying io.WriteCloser, ensuring no other
// MultiWriteCloserFanIn can interrupt this one's writing.
func (fanin *MultiWriteCloserFanIn) Write(data []byte) (int, error) {
	fanin.pLock.Lock()
	var err error
	var written, m int
	for err == nil && written < len(data) {
		m, err = fanin.iowc.Write(data[written:])
		written += m
	}
	fanin.pLock.Unlock()
	return written, err
}

// Close marks the MultiWriteCloserFanIn as finished.  The last Close method invoked for a group of
// MultiWriteCloserFanIn instances will trigger a close of the underlying io.WriteCloser.
func (fanin *MultiWriteCloserFanIn) Close() error {
	fanin.done.Done()
	return nil
}
