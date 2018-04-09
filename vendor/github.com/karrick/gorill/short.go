package gorill

import "io"

// ShortReadWriteCloser wraps a io.ReadWriteCloser, but the Read and Write operations cannot exceed
// the MaxRead and MaxWrite sizes.
type ShortReadWriteCloser struct {
	io.Reader
	io.WriteCloser
	MaxRead  int
	MaxWrite int
}

// Read reads from the wrapped io.Reader, but returns EOF if attempts to read beyond the MaxRead.
func (s ShortReadWriteCloser) Read(buf []byte) (int, error) {
	var short bool
	index := len(buf)
	if index > s.MaxRead {
		index = s.MaxRead
		short = true
	}
	n, err := s.Reader.Read(buf[:index])
	if short {
		return n, io.ErrUnexpectedEOF
	}
	return n, err
}

// Write writes from the wrapped io.Writer, but returns EOF if attempts to write beyond the MaxWrite.
func (s ShortReadWriteCloser) Write(data []byte) (int, error) {
	var short bool
	index := len(data)
	if index > s.MaxWrite {
		index = s.MaxWrite
		short = true
	}
	n, err := s.WriteCloser.Write(data[:index])
	if short {
		return n, io.ErrShortWrite
	}
	return n, err
}

// ShortReadWriter wraps a io.Reader and io.Writer, but the Read and Write operations cannot exceed
// the MaxRead and MaxWrite sizes.
type ShortReadWriter struct {
	io.Reader
	io.Writer
	MaxRead  int
	MaxWrite int
}

// Read reads from the wrapped io.Reader, but returns EOF if attempts to read beyond the MaxRead.
func (s ShortReadWriter) Read(buf []byte) (int, error) {
	var short bool
	index := len(buf)
	if index > s.MaxRead {
		index = s.MaxRead
		short = true
	}
	n, err := s.Reader.Read(buf[:index])
	if short {
		return n, io.ErrUnexpectedEOF
	}
	return n, err
}

// Write writes from the wrapped io.Writer, but returns EOF if attempts to write beyond the MaxWrite.
func (s ShortReadWriter) Write(data []byte) (int, error) {
	var short bool
	index := len(data)
	if index > s.MaxWrite {
		index = s.MaxWrite
		short = true
	}
	n, err := s.Writer.Write(data[:index])
	if short {
		return n, io.ErrShortWrite
	}
	return n, err
}

// ShortWriter returns a structure that wraps an io.Writer, but returns io.ErrShortWrite when the
// number of bytes to write exceeds a preset limit.
//
//   bb := gorill.NopCloseBuffer()
//   sw := gorill.ShortWriter(bb, 16)
//
//   n, err := sw.Write([]byte("short write"))
//   // n == 11, err == nil
//
//   n, err := sw.Write([]byte("a somewhat longer write"))
//   // n == 16, err == io.ErrShortWrite
func ShortWriter(w io.Writer, max int) io.Writer {
	return shortWriter{Writer: w, max: max}
}

func (s shortWriter) Write(data []byte) (int, error) {
	var short bool
	index := len(data)
	if index > s.max {
		index = s.max
		short = true
	}
	n, err := s.Writer.Write(data[:index])
	if short {
		return n, io.ErrShortWrite
	}
	return n, err
}

type shortWriter struct {
	io.Writer
	max int
}

// ShortWriteCloser returns a structure that wraps an io.WriteCloser, but returns io.ErrShortWrite
// when the number of bytes to write exceeds a preset limit.
//
//   bb := gorill.NopCloseBuffer()
//   sw := gorill.ShortWriteCloser(bb, 16)
//
//   n, err := sw.Write([]byte("short write"))
//   // n == 11, err == nil
//
//   n, err := sw.Write([]byte("a somewhat longer write"))
//   // n == 16, err == io.ErrShortWrite
func ShortWriteCloser(iowc io.WriteCloser, max int) io.WriteCloser {
	return shortWriteCloser{WriteCloser: iowc, max: max}
}

func (s shortWriteCloser) Write(data []byte) (int, error) {
	var short bool
	index := len(data)
	if index > s.max {
		index = s.max
		short = true
	}
	n, err := s.WriteCloser.Write(data[:index])
	if short {
		return n, io.ErrShortWrite
	}
	return n, err
}

type shortWriteCloser struct {
	io.WriteCloser
	max int
}
