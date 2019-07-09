package gorill

import (
	"io"
	"io/ioutil"
)

// ReadAllThenClose reads all bytes from rc then closes it.  It returns any
// errors that occurred when either reading or closing rc.
func ReadAllThenClose(rc io.ReadCloser) ([]byte, error) {
	buf, rerr := ioutil.ReadAll(rc)
	cerr := rc.Close() // always close regardless of read error
	if rerr != nil {
		return nil, rerr // Read error has more context than Close error
	}
	if cerr != nil {
		return nil, cerr
	}
	return buf, nil
}
