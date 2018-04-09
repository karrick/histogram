package gorill

import (
	"bytes"
	"testing"
	"time"
)

func TestSpooledWriteCloserFlushForcesBytesWritten(t *testing.T) {
	test := func(buf []byte, flushPeriodicity time.Duration) {
		bb := new(bytes.Buffer)

		SlowWriter := SlowWriter(bb, 10*time.Millisecond)
		spoolWriter, _ := NewSpooledWriteCloser(NopCloseWriter(SlowWriter), Flush(flushPeriodicity))
		defer func() {
			if err := spoolWriter.Close(); err != nil {
				t.Errorf("Actual: %s; Expected: %#v", err, nil)
			}
		}()

		n, err := spoolWriter.Write(buf)
		if want := len(buf); n != want {
			t.Errorf("Actual: %#v; Expected: %#v", n, want)
		}
		if err != nil {
			t.Errorf("Actual: %#v; Expected: %#v", err, nil)
		}
		if err = spoolWriter.Flush(); err != nil {
			t.Errorf("Actual: %s; Expected: %#v", err, nil)
		}
		if want := string(buf); bb.String() != want {
			t.Errorf("Actual: %#v; Expected: %#v", bb.String(), want)
		}
	}
	test(smallBuf, time.Millisecond)
	test(largeBuf, time.Millisecond)

	test(smallBuf, time.Hour)
	test(largeBuf, time.Hour)
}

func TestSpooledWriteCloserCloseCausesFlush(t *testing.T) {
	test := func(buf []byte, flushPeriodicity time.Duration) {
		bb := NewNopCloseBuffer()

		spoolWriter, _ := NewSpooledWriteCloser(bb, Flush(flushPeriodicity))

		n, err := spoolWriter.Write(buf)
		if want := len(buf); n != want {
			t.Errorf("Actual: %#v; Expected: %#v", n, want)
		}
		if err != nil {
			t.Errorf("Actual: %#v; Expected: %#v", err, nil)
		}
		if err := spoolWriter.Close(); err != nil {
			t.Errorf("Actual: %s; Expected: %#v", err, nil)
		}
		if want := string(buf); bb.String() != want {
			t.Errorf("Actual: %#v; Expected: %#v", bb.String(), want)
		}
	}
	test(smallBuf, time.Millisecond)
	test(largeBuf, time.Millisecond)

	test(smallBuf, time.Hour)
	test(largeBuf, time.Hour)
}
