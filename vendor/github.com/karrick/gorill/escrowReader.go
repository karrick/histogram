package gorill

import (
	"bytes"
	"io"
)

// EscrowReader is a structure that mimics the io.ReadCloser interface, yet
// already has all the payload bytes stored in memory along with any read error
// that took place while reading the payload. The benefit of using it is that
// other code may re-read the payload without an additional penalty, as the
// bytes are already buffered in memory.
type EscrowReader struct {
	// off holds the offset into the buffer that Read will return bytes from.
	off int64

	// cerr holds any error that occurred whiel closing the data source.
	cerr error

	// rerr holds any error that occurred while reading the data source.
	rerr error

	buf []byte // payload holds the request body payload.
}

// NewEscrowReader reads and consumes all the data from the specified
// io.ReadCloser into either a new bytes.Buffer or a specified bytes.Buffer,
// then returns an io.ReadCloser that allows the data to be read multiple
// times. It always closes the provided io.ReadCloser.
//
// It does not return any errors during instantiation, because any read error
// encountered will be returned after the last byte is read from the provided
// io.ReadCloser. Likewise any close error will be returned by the structure's
// Close method.
//
//     func someHandler(w http.ResponseWriter, r *http.Request) {
//         // Get a scratch buffer for the example. For production code, consider using
//         // a free-list of buffers, such as https://github.com/karrick/gobp
//         bb := new(bytes.Buffer)
//         r.Body = NewEscrowReader(r.Body, bb)
//         // ...
//     }
func NewEscrowReader(iorc io.ReadCloser, bb *bytes.Buffer) *EscrowReader {
	if bb == nil {
		bb = new(bytes.Buffer)
	}
	_, rerr := bb.ReadFrom(iorc)
	if rerr == nil {
		// Mimic expected behavior of returning io.EOF when there are no bytes
		// remaining to be read.
		rerr = io.EOF
	}
	cerr := iorc.Close()
	return &EscrowReader{buf: bb.Bytes(), cerr: cerr, rerr: rerr}
}

// Bytes returns the slice of bytes read from the original data source.
func (er *EscrowReader) Bytes() []byte { return er.buf }

// Close returns the error that took place when closing the original
// io.ReadCloser. Under normal circumstances it will be nil.
func (er *EscrowReader) Close() error { return er.cerr }

// Err returns the error encountered while reading from the source io.ReadCloser
// if not io.EOF; otherwise it returns the error encountered while closing it.
// This method comes in handy when you know you have an EscrowReader, and you
// want to know whether the entire payload was slurped in.
//
//     func example(iorc io.ReadCloser) ([]byte, error) {
//         if er, ok := iorc.(*gorill.EscrowReader); ok {
//             return er.Bytes(), er.Err()
//         }
//     }
func (er *EscrowReader) Err() error {
	// An error encountered while reading has more context than an error
	// encountered while closing.
	if er.rerr != io.EOF {
		return er.rerr
	}
	return er.cerr
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0
// <= n <= len(p)) and any error encountered. Even if Read returns n < len(p),
// it may use all of p as scratch space during the call.  If some data is
// available but not len(p) bytes, Read conventionally returns what is available
// instead of waiting for more.
//
// When there are no more bytes to be read from the buffer, it will return any
// error encountered when reading the original io.ReadCloser data source. That
// error value is normally nil, but could be any other error other than io.EOF.
func (er *EscrowReader) Read(b []byte) (int, error) {
	if er.off >= int64(len(er.buf)) {
		// Once everything has been read, any further reads return the error
		// that was recorded when slurping in the original data source.
		return 0, er.rerr
	}
	n := copy(b, er.buf[er.off:])
	er.off += int64(n)
	return n, nil
}

// Reset will cause the next Read to read from the beginning of the buffer.
func (er *EscrowReader) Reset() {
	er.off = 0
}

// WriteTo writes the entire buffer contents to w, and returns the number of
// bytes written along with any error.
func (er *EscrowReader) WriteTo(w io.Writer) (int64, error) {
	var n int64
	if len(er.buf) > 0 {
		nw, err := w.Write(er.buf)
		n = int64(nw)
		if err == nil && nw != len(er.buf) {
			// While io.Writer function is supposed to write all the bytes or
			// return an error describing why, protect against a misbehaving
			// writer that writes fewer bytes than requested and returns no
			// error.
			return n, io.ErrShortWrite
		}
		return n, err
	}
	return n, nil
}
