package gorill

import "strings"

// ErrList is a slice of errors, useful when a function must return a single error, but has multiple
// independent errors to return.
type ErrList []error

// Append appends non-nil errors to the list of errors.
func (e *ErrList) Append(b error) {
	if b != nil {
		*e = append(*e, b)
	}
}

// Count returns number of non-nil errors accumulated in ErrList.
func (e ErrList) Count() int {
	return len([]error(e))
}

// Err returns either a list of non-nil error values, or a single error value if the list only
// contains one error.
func (e ErrList) Err() error {
	errors := make([]error, 0, len([]error(e)))
	for _, e := range []error(e) {
		if e != nil {
			errors = append(errors, e)
		}
	}
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return e[0]
	default:
		return ErrList(errors)
	}
}

// Error returns the string version of an error list, which is the list of errors, joined by a
// comma-space byte sequence.
func (e ErrList) Error() string {
	es := make([]string, 0, len([]error(e)))
	for i := range []error(e) {
		if []error(e)[i] != nil {
			es = append(es, []error(e)[i].Error())
		}
	}
	return strings.Join(es, ", ")
}
