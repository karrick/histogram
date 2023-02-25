package gorill

// newline returns a string with exactly one terminating newline character.
// More simple than strings.TrimRight.  When input string ends with multiple
// newline characters, it will strip off all but first one, reusing the same
// underlying string bytes.  When string does not end in a newline character, it
// returns the original string with a newline character appended.  Newline
// characters before any non-newline characters are ignored.
func newline(s string) string {
	l := len(s)
	if l == 0 {
		return "\n"
	}

	// While this is O(length s), it stops as soon as it finds the first non
	// newline character in the string starting from the right hand side of the
	// input string.  Generally this only scans one or two characters and
	// returns.

	for i := l - 1; i >= 0; i-- {
		if s[i] != '\n' {
			if i+1 < l && s[i+1] == '\n' {
				return s[:i+2]
			}
			return s[:i+1] + "\n"
		}
	}

	return s[:1] // all newline characters, so just return the first one
}
