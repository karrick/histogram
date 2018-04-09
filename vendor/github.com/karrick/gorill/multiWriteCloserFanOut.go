package gorill

import (
	"io"
	"sync"
)

// MultiWriteCloserFanOut is a structure that allows additions to and removals from the list of
// io.WriteCloser objects that will be written to.
type MultiWriteCloserFanOut struct {
	lock        sync.RWMutex
	writerMap   map[io.WriteCloser]struct{}
	writerSlice []io.WriteCloser
}

// NewMultiWriteCloserFanOut returns a MultiWriteCloserFanOut that is go-routine safe.
//
//   bb1 = gorill.NewNopCloseBuffer()
//   bb2 = gorill.NewNopCloseBuffer()
//   mw = gorill.NewMultiWriteCloserFanOut(bb1, bb2)
//   n, err := mw.Write([]byte("blob"))
//   if want := 4; n != want {
//   	t.Errorf("Actual: %#v; Expected: %#v", n, want)
//   }
//   if err != nil {
//   	t.Errorf("Actual: %#v; Expected: %#v", err, nil)
//   }
//   if want := "blob"; bb1.String() != want {
//   	t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
//   }
//   if want := "blob"; bb2.String() != want {
//   	t.Errorf("Actual: %#v; Expected: %#v", bb2.String(), want)
//   }
func NewMultiWriteCloserFanOut(writers ...io.WriteCloser) *MultiWriteCloserFanOut {
	mwc := &MultiWriteCloserFanOut{writerMap: make(map[io.WriteCloser]struct{})}
	for _, w := range writers {
		mwc.writerMap[w] = struct{}{}
	}
	mwc.update()
	return mwc
}

// update makes the slice reflect contents of the altered map
func (mwc *MultiWriteCloserFanOut) update() {
	mwc.writerSlice = make([]io.WriteCloser, 0, len(mwc.writerMap))
	for iow := range mwc.writerMap {
		mwc.writerSlice = append(mwc.writerSlice, iow)
	}
}

// Add adds an io.WriteCloser to the list of writers to be written to whenever this
// MultiWriteCloserFanOut is written to.  It returns the number of io.WriteCloser instances attached
// to the MultiWriteCloserFanOut instance.
//
//   bb1 = gorill.NewNopCloseBuffer()
//   mw = gorill.NewMultiWriteCloserFanOut(bb1)
//   bb2 = gorill.NewNopCloseBuffer()
//   mw.Add(bb2)
func (mwc *MultiWriteCloserFanOut) Add(w io.WriteCloser) int {
	mwc.lock.Lock()
	defer mwc.lock.Unlock()

	mwc.writerMap[w] = struct{}{}
	mwc.update()
	return len(mwc.writerSlice)
}

// Close will close the underlying io.WriteCloser, and releases resources.
func (mwc *MultiWriteCloserFanOut) Close() error {
	mwc.lock.Lock()
	defer mwc.lock.Unlock()

	var errors ErrList
	for _, iowc := range mwc.writerSlice {
		errors.Append(iowc.Close())
	}
	return errors.Err()
}

// Count returns the number of io.WriteCloser instances attached to the MultiWriteCloserFanOut
// instance.
//
//   mw = gorill.NewMultiWriteCloserFanOut()
//   count := mw.Count() // returns 1
//   mw.Add(gorill.NewNopCloseBuffer())
//   count = mw.Count() // returns 2
func (mwc *MultiWriteCloserFanOut) Count() int {
	mwc.lock.RLock()
	defer mwc.lock.RUnlock()

	return len(mwc.writerSlice)
}

// IsEmpty returns true if and only if there are no writers in the list of writers to be written to.
//
//   mw = gorill.NewMultiWriteCloserFanOut()
//   mw.IsEmpty() // returns true
//   mw.Add(gorill.NewNopCloseBuffer())
//   mw.IsEmpty() // returns false
func (mwc *MultiWriteCloserFanOut) IsEmpty() bool {
	mwc.lock.RLock()
	defer mwc.lock.RUnlock()

	return len(mwc.writerSlice) == 0
}

// Remove removes an io.WriteCloser from the list of writers to be written to whenever this
// MultiWriteCloserFanOut is written to.  It returns the number of io.WriteCloser instances attached
// to the MultiWriteCloserFanOut instance.
//
//   bb1 = gorill.NewNopCloseBuffer()
//   bb2 = gorill.NewNopCloseBuffer()
//   mw = gorill.NewMultiWriteCloserFanOut(bb1, bb2)
//   remaining := mw.Remove(bb1) // returns 1
func (mwc *MultiWriteCloserFanOut) Remove(w io.WriteCloser) int {
	mwc.lock.Lock()
	defer mwc.lock.Unlock()

	delete(mwc.writerMap, w)
	mwc.update()
	return len(mwc.writerSlice)
}

// Write writes the data to all the writers in the MultiWriteCloserFanOut.  It removes and invokes Close
// method for all io.WriteClosers that returns an error when written to.
//
//   bb1 = gorill.NewNopCloseBuffer()
//   bb2 = gorill.NewNopCloseBuffer()
//   mw = gorill.NewMultiWriteCloserFanOut(bb1, bb2)
//   n, err := mw.Write([]byte("blob"))
//   if want := 4; n != want {
//   	t.Errorf("Actual: %#v; Expected: %#v", n, want)
//   }
//   if err != nil {
//   	t.Errorf("Actual: %#v; Expected: %#v", err, nil)
//   }
func (mwc *MultiWriteCloserFanOut) Write(data []byte) (int, error) {
	mwc.lock.RLock()
	defer mwc.lock.RUnlock()

	// NOTE: the complexity of wait group and go routines does not
	// solve the slow writer problem, but it helps
	var lock sync.Mutex
	var wg sync.WaitGroup
	var errored []io.WriteCloser
	wg.Add(len(mwc.writerSlice))
	for _, sw := range mwc.writerSlice {
		go func(w io.WriteCloser) {
			n, err := w.Write(data)
			if n != len(data) {
				err = io.ErrShortWrite
			}
			if err != nil {
				lock.Lock()
				errored = append(errored, w)
				lock.Unlock()
			}
			wg.Done()
		}(sw)
	}
	wg.Wait()
	if len(errored) > 0 {
		for _, w := range errored {
			delete(mwc.writerMap, w)
			w.Close() // BUG might cause bug when client tries to later Close ???
		}
		mwc.update()
	}
	return len(data), nil
}
