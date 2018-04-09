package gorill

type opcode byte

const (
	_read opcode = iota
	_write
	_flush
)

// rillJob represents a job to perform either a read or write operation to a stream
type rillJob struct {
	op      opcode
	data    []byte
	results chan rillResult
}

func newRillJob(op opcode, data []byte) *rillJob {
	return &rillJob{op: op, data: data, results: make(chan rillResult, 1)}
}

// rillResult represents the return values for a read or write operation to a stream
type rillResult struct {
	n   int
	err error
}

// ErrReadAfterClose is returned if a Read is attempted after Close called.
type ErrReadAfterClose struct{}

// Error returns a string representation of a ErrReadAfterClose error instance.
func (e ErrReadAfterClose) Error() string {
	return "read on closed reader"
}

// ErrWriteAfterClose is returned if a Write is attempted after Close called.
type ErrWriteAfterClose struct{}

// Error returns a string representation of a ErrWriteAfterClose error instance.
func (e ErrWriteAfterClose) Error() string {
	return "write on closed writer"
}
