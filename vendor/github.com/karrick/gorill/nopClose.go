package gorill

import (
	"bytes"
	"io"
)

// NewNopCloseBuffer returns a structure that wraps bytes.Buffer with a no-op Close method.  It can
// be used in tests that need a bytes.Buffer, but need to provide a Close method.
//
//   bb := gorill.NopCloseBuffer()
//   bb.Write([]byte("example"))
//   bb.Close() // does nothing
func NewNopCloseBuffer() *NopCloseBuffer {
	return &NopCloseBuffer{Buffer: new(bytes.Buffer), closed: false}
}

// NewNopCloseBufferSize returns a structure that wraps bytes.Buffer with a no-op Close method,
// using a specified buffer size.  It can be used in tests that need a bytes.Buffer, but need to
// provide a Close method.
//
//   bb := gorill.NopCloseBufferSize(8192)
//   bb.Write([]byte("example"))
//   bb.Close() // does nothing
func NewNopCloseBufferSize(size int) *NopCloseBuffer {
	return &NopCloseBuffer{Buffer: bytes.NewBuffer(make([]byte, 0, size)), closed: false}
}

// NopCloseBuffer is a structure that wraps a buffer, but also provides a no-op Close method.
type NopCloseBuffer struct {
	*bytes.Buffer
	closed bool
}

// Close returns nil error.
func (m *NopCloseBuffer) Close() error { m.closed = true; return nil }

// IsClosed returns false, unless NopCloseBuffer's Close method has been invoked
func (m *NopCloseBuffer) IsClosed() bool { return m.closed }

// NopCloseReader returns a structure that implements io.ReadCloser, but provides a no-op Close
// method.  It is useful when you have an io.Reader that you must pass to a method that requires an
// io.ReadCloser.  It is the same as ioutil.NopCloser, but for provided here for symmetry with
// NopCloseWriter.
//
//   iorc := gorill.NopCloseReader(ior)
//   iorc.Close() // does nothing
func NopCloseReader(ior io.Reader) io.ReadCloser { return nopCloseReader{ior} }

type nopCloseReader struct{ io.Reader }

func (nopCloseReader) Close() error { return nil }

// NopCloseWriter returns a structure that implements io.WriteCloser, but provides a no-op Close
// method.  It is useful when you have an io.Writer that you must pass to a method that requires an
// io.WriteCloser.  It is the counter-part to ioutil.NopCloser, but for io.Writer.
//
//   iowc := gorill.NopCloseWriter(iow)
//   iowc.Close() // does nothing
func NopCloseWriter(iow io.Writer) io.WriteCloser { return nopCloseWriter{iow} }

func (nopCloseWriter) Close() error { return nil }

type nopCloseWriter struct{ io.Writer }
