package gorill

import (
	"io"
	"os"
)

// FilesReader is an io.ReadCloser that can be used to read over the contents of
// all of the files specified by pathnames. It only opens a single file handle
// at a time. When reading from the currently open file handle returns io.EOF,
// it closes that file handle, and the next Read will cause the following file
// in the series to be opened and read from.
type FilesReader struct {
	// Pathnames is a list of remaining files to read. It is kept up to date
	// where when a file is opened, its name is removed from the head of the
	// list.
	Pathnames []string

	// fh is the currently open file handle from which Read operations will take
	// place. It can be nil, in which case Read will attempt to open the next
	// file in the series and read from it.
	fh *os.File
}

// Close forgets the list of remaining files in the series, then closes the
// currently open file handle, returning any error from the operating system.
func (fr *FilesReader) Close() error {
	fr.Pathnames = nil
	if fr.fh == nil {
		return nil
	}
	err := fr.fh.Close()
	fr.fh = nil
	return err
}

// Read reads up to len(p) bytes into p. It returns the number of bytes read (0
// <= n <= len(p) and any error encountered.
func (fr *FilesReader) Read(p []byte) (int, error) {
	if fr.fh == nil {
		if err := fr.next(); err != nil {
			return 0, err
		}
	}
	for {
		nr, err := fr.fh.Read(p)
		if err == io.EOF {
			if err = fr.fh.Close(); err != nil {
				return nr, err
			}
			if err = fr.next(); err != nil {
				return nr, err
			}
			if nr == 0 {
				continue
			}
		}
		return nr, err
	}
}

// Next closes the currently open file handle and opens the next file in the
// series. If there are no files left it returns io.EOF. It can be used to skip
// the remaining contents of the currently open file. Additional Read operations
// will be invoked against the following file in the series, if non empty.
func (fr *FilesReader) Next() error {
	if fr.fh != nil {
		err := fr.fh.Close()
		fr.fh = nil
		if err != nil {
			return err
		}
	}
	return fr.next()
}

// next opens the next file in the series. If there are no files left it returns
// io.EOF.
func (fr *FilesReader) next() error {
	if len(fr.Pathnames) == 0 {
		return io.EOF
	}

	// NOTE: Consider having the '-' string open standard input
	fh, err := os.Open(fr.Pathnames[0])
	if err != nil {
		return err
	}

	fr.fh = fh
	fr.Pathnames = fr.Pathnames[1:]
	return nil
}
