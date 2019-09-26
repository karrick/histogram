package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/karrick/golf"
)

var (
	ProgramName            string
	ProgramLongDescription string
	ProgramOneLineSummary  string
	ProgramUsageExamples   string
	ProgramVersion         = "0.0.0"

	optHelp    = golf.BoolP('h', "help", false, "Print command line help and exit")
	optVersion = golf.BoolP('V', "version", false, "Print version information and exit")

	optQuiet   = golf.BoolP('q', "quiet", false, "Do not print intermediate errors to stderr")
	optVerbose = golf.BoolP('v', "verbose", false, "Print verbose output to stderr")
)

func briefUsage() {
	fmt.Fprintf(os.Stderr, "Use `%s --help` for more information.\n", ProgramName)
}

func init() {
	var err error
	ProgramName, err = os.Executable()
	if err != nil {
		ProgramName = filepath.Base(os.Args[0])
	}

	// Rather than display the entire usage information for a parsing error,
	// merely allow golf library to display the error message, then print the
	// command the user may use to show command line usage information.
	golf.Usage = briefUsage
}

func main() {
	golf.Parse()

	if *optHelp {
		fmt.Fprintf(os.Stderr, "%s version %s\n\n", ProgramName, ProgramVersion)
		if ProgramLongDescription != "" {
			fmt.Fprintln(os.Stderr, ProgramLongDescription)
		}
		if ProgramUsageExamples != "" {
			fmt.Fprintln(os.Stderr, ProgramUsageExamples)
		}
		fmt.Fprintln(os.Stderr, "Command line options:")
		golf.PrintDefaults()
		os.Exit(0)
	}

	if *optVersion {
		fmt.Fprintf(os.Stderr, "%s version %s; %s; `%s --help` for more information.\n", ProgramName, ProgramVersion, ProgramOneLineSummary, ProgramName)
		os.Exit(0)
	}

	if err := cmd(); err != nil {
		if _, ok := err.(ErrUsage); ok {
			fmt.Fprintf(os.Stderr, "%s: %s\n", ProgramName, err)
			briefUsage()
			os.Exit(2)
		}
		stderr("%s", newline(err.Error()))
		os.Exit(1)
	}
}

// newline returns a string with exactly one terminating newline character.
// More simple than strings.TrimRight.  When input string has multiple newline
// characters, it will strip off all but first one, reusing the same underlying
// string bytes.  When string does not end in a newline character, it returns
// the original string with a newline character appended.
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

// stderr formats and prints its arguments to standard error after prefixing
// them with the program name.
func stderr(f string, args ...interface{}) {
	os.Stderr.Write([]byte(ProgramName + ": " + fmt.Sprintf(f, args...)))
}

// verbose formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optVerbose is false.
func verbose(f string, args ...interface{}) {
	if *optVerbose {
		os.Stderr.Write([]byte(ProgramName + ": " + fmt.Sprintf(f, args...)))
	}
}

// warning formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optQuiet is true.
func warning(f string, args ...interface{}) {
	if !*optQuiet {
		os.Stderr.Write([]byte(ProgramName + ": " + fmt.Sprintf(f, args...)))
	}
}

// ErrUsage is returned by the program code it discovers a usage error when
// parsing command line arguments.  It displays the error message, prints the
// program usage information, then exists with program status code set to 2.
type ErrUsage string

func (e ErrUsage) Error() string { return string(e) }

func NewErrUsage(f string, a ...interface{}) ErrUsage {
	return ErrUsage(fmt.Sprintf(f, a...))
}
