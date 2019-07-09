package gorill

import (
	"io"
)

// LineTerminatedReader reads from the source io.Reader and ensures the final
// byte from this is a newline.
type LineTerminatedReader struct {
	R                   io.Reader
	wasFinalByteNewline bool
	oweNewline          bool
}

// Read satisfies the io.Reader interface by reading up to len(p) bytes into p.
// It returns the number of bytes read (0 <= n <= len(p)) and any error
// encountered.
func (r *LineTerminatedReader) Read(p []byte) (int, error) {
	if r.oweNewline {
		r.oweNewline = false
		if len(p) == 0 {
			return 0, nil // from io.Reader documentation
		}
		p[0] = '\n'
		return 1, io.EOF // allowed per io.Reader documentation
	}
	n, err := r.R.Read(p)
	if n > 0 {
		r.wasFinalByteNewline = p[n-1] == '\n'
	}
	if err != io.EOF || r.wasFinalByteNewline {
		return n, err
	}
	// fmt.Fprintf(os.Stderr, "Hit EOF && final byte was not a newline.\n")
	if n == len(p) {
		// No room to append newline to this buffer.
		r.oweNewline = true
		return n, nil
	}
	// Append the newline byte when it fits.
	p[n] = '\n'
	return n + 1, err
}
