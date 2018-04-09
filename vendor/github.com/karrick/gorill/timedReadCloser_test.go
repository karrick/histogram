package gorill

import (
	"bytes"
	"testing"
	"time"
)

func TestTimedReadCloser(t *testing.T) {
	corpus := "this is a test"
	bb := bytes.NewReader([]byte(corpus))
	rc := NewTimedReadCloser(NopCloseReader(bb), time.Second)
	defer rc.Close()

	buf := make([]byte, 1000)
	n, err := rc.Read(buf)
	if actual, want := n, len(corpus); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if actual, want := string(buf[:n]), corpus; actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if actual, want := err, error(nil); actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
}

func TestTimedReadCloserTimesOut(t *testing.T) {
	corpus := "this is a test"
	bb := bytes.NewReader([]byte(corpus))
	sr := SlowReader(bb, 10*time.Millisecond)
	rc := NewTimedReadCloser(NopCloseReader(sr), time.Millisecond)
	defer rc.Close()

	buf := make([]byte, 1000)
	n, err := rc.Read(buf)
	if actual, want := n, 0; actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if actual, want := string(buf[:n]), ""; actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if actual, want := err, ErrTimeout(time.Millisecond); actual != want {
		t.Errorf("Actual: %s; Expected: %s", actual, want)
	}
}

func TestTimedReadCloserReadAfterCloseReturnsError(t *testing.T) {
	bb := NewNopCloseBuffer()
	rc := NewTimedReadCloser(NopCloseReader(bb), time.Millisecond)
	rc.Close()

	buf := make([]byte, 1000)
	n, err := rc.Read(buf)
	if actual, want := n, 0; actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if actual, want := string(buf[:n]), ""; actual != want {
		t.Errorf("Actual: %#v; Expected: %#v", actual, want)
	}
	if _, ok := err.(ErrReadAfterClose); err == nil || !ok {
		t.Errorf("Actual: %s; Expected: %#v", err, ErrReadAfterClose{})
	}
}
