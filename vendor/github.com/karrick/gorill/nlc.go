package gorill

import (
	"bytes"
	"io"
)

// NewlineCounter counts the number of lines from the io.Reader, returning the
// same number of lines read regardless of whether the final Read terminated in
// a newline character.
func NewlineCounter(ior io.Reader) (int, error) {
	var newlines, total, n int
	var isNotFinalNewline bool
	var err error
	var reserved [4096]byte // allocate buffer space on the call stack
	buf := reserved[:]      // create slice using pre-allocated array from reserved

	for {
		n, err = ior.Read(buf)
		if n > 0 {
			total += n
			isNotFinalNewline = buf[n-1] != '\n'
			var searchOffset int
			for {
				index := bytes.IndexByte(buf[searchOffset:n], '\n')
				if index == -1 {
					break // done counting newlines from this chunk
				}
				newlines++                // count this newline
				searchOffset += index + 1 // start next search following this newline
			}
		}
		if err != nil {
			if err == io.EOF {
				err = nil // io.EOF is expected at end of stream
			}
			break // do not try to read more if error
		}
	}

	// Return the same number of lines read regardless of whether the final read
	// terminated in a newline character.
	if isNotFinalNewline {
		newlines++
	} else if total == 1 {
		newlines--
	}
	return newlines, err
}
