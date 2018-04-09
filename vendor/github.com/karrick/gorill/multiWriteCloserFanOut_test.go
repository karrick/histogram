package gorill

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestMultiWriteCloserFanOutNoWriteClosers(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()
	if want := 0; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	mw.Close()
}

func TestMultiWriteCloserFanOutNewWriteCloser(t *testing.T) {
	bb1 := NewNopCloseBuffer()
	bb2 := NewNopCloseBuffer()
	mw := NewMultiWriteCloserFanOut(bb1, bb2)

	if want := 2; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := "blob"; bb1.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
	}
	if want := "blob"; bb2.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb2.String(), want)
	}
	mw.Close()
}

func TestMultiWriteCloserFanOutOneWriteCloser(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()

	bb1 := NewNopCloseBuffer()
	if actual, expected := mw.Add(bb1), 1; actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	if want := 1; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := "blob"; bb1.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
	}
	mw.Close()
}

func TestMultiWriteCloserFanOutTwoWriteClosers(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()

	bb1 := NewNopCloseBuffer()
	mw.Add(bb1)
	bb2 := NewNopCloseBuffer()
	mw.Add(bb2)
	if want := 2; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := "blob"; bb1.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
	}
	if want := "blob"; bb2.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb2.String(), want)
	}
	mw.Close()
}

func TestMultiWriteCloserFanOutRemoveWriteCloser(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()

	bb1 := NewNopCloseBuffer()
	mw.Add(bb1)
	bb2 := NewNopCloseBuffer()
	mw.Add(bb2)
	if actual, expected := mw.Remove(bb1), 1; actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := ""; bb1.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
	}
	if want := "blob"; bb2.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb2.String(), want)
	}
	// remove last one, should return true
	if actual, expected := mw.Remove(bb2), 0; actual != expected {
		t.Errorf("Actual: %#v; Expected: %#v", actual, expected)
	}
	mw.Close()
}

func TestMultiWriteCloserFanOutRemoveEveryWriteCloser(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()

	bb1 := NewNopCloseBuffer()
	mw.Add(bb1)
	bb2 := NewNopCloseBuffer()
	mw.Add(bb2)
	mw.Remove(bb1)
	mw.Remove(bb2)
	if want := 0; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	n, err := mw.Write([]byte("blob"))
	if want := 4; n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := ""; bb1.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb1.String(), want)
	}
	if want := ""; bb2.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", bb2.String(), want)
	}
	mw.Close()
}

type testWriteCloser struct {
	closed bool
}

func (w *testWriteCloser) Write([]byte) (int, error) {
	return 0, io.ErrShortWrite
}

func (w *testWriteCloser) Close() error {
	w.closed = true
	return nil
}

func (w *testWriteCloser) IsClosed() bool {
	return w.closed
}

func TestMultiWriteCloserFanOutWriteErrorRemovesBadWriteCloser(t *testing.T) {
	mw := NewMultiWriteCloserFanOut()

	buf := NewNopCloseBuffer()
	ew := &testWriteCloser{}

	mw.Add(buf)
	mw.Add(ew)

	n, err := mw.Write([]byte(alphabet))
	if want := len(alphabet); n != want {
		t.Errorf("Actual: %#v; Expected: %#v", n, want)
	}
	if err != nil {
		t.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	if want := alphabet; buf.String() != want {
		t.Errorf("Actual: %#v; Expected: %#v", buf.String(), want)
	}

	mw.Remove(buf)
	// NOTE: testWriteCloser should have been removed during error write
	if want := true; ew.IsClosed() != want {
		t.Errorf("Actual: %#v; Expected: %#v", ew.IsClosed(), want)
	}
	// NOTE: testWriteCloser should have been removed during error write
	if want := 0; mw.Count() != want {
		t.Errorf("Actual: %#v; Expected: %#v", mw.Count(), want)
	}
	mw.Close()
}

const writersCount = 1000
const data = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"

func BenchmarkWriterMultiWriteCloserFanOutWrite(b *testing.B) {
	mw := NewMultiWriteCloserFanOut()
	for i := 0; i < writersCount; i++ {
		mw.Add(NewNopCloseBuffer())
	}
	n, err := mw.Write([]byte(data))
	if n != len(data) {
		b.Errorf("Actual: %#v; Expected: %#v", n, 4)
	}
	if err != nil {
		b.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	mw.Close()
}

func BenchmarkWriterMultiWriteCloserFanOutWriteSlow(b *testing.B) {
	mw := NewMultiWriteCloserFanOut()
	for i := 0; i < writersCount; i++ {
		mw.Add(NopCloseWriter(SlowWriter(new(bytes.Buffer), 10*time.Millisecond)))
	}
	n, err := mw.Write([]byte(data))
	if n != len(data) {
		b.Errorf("Actual: %#v; Expected: %#v", n, 4)
	}
	if err != nil {
		b.Errorf("Actual: %#v; Expected: %#v", err, nil)
	}
	mw.Close()
}
