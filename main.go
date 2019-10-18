package main // import "github.com/karrick/histogram"

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/gobls"
	"github.com/karrick/gohistogram"
	"github.com/karrick/golf"
	"github.com/karrick/gorill"
	"github.com/karrick/gows"
)

// fatal prints the error to standard error then exits the program with status
// code 1.
func fatal(err error) {
	stderr("%s\n", err)
	os.Exit(1)
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
	os.Stderr.Write([]byte(ProgramName + ": " + newline(fmt.Sprintf(f, args...))))
}

// usage prints the error to standard error, prints message how to get help,
// then exits the program with status code 2.
func usage(f string, args ...interface{}) {
	stderr(f, args...)
	golf.Usage()
	os.Exit(2)
}

// verbose formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optVerbose is false.
func verbose(f string, args ...interface{}) {
	if *optVerbose {
		stderr(f, args...)
	}
}

// warning formats and prints its arguments to standard error after prefixing
// them with the program name.  This skips printing when optQuiet is true.
func warning(f string, args ...interface{}) {
	if !*optQuiet {
		stderr(f, args...)
	}
}

var ProgramName string

func init() {
	var err error
	if ProgramName, err = os.Executable(); err != nil {
		ProgramName = os.Args[0]
	}
	ProgramName = filepath.Base(ProgramName)

	// Rather than display the entire usage information for a parsing error,
	// merely allow golf library to display the error message, then print the
	// command the user may use to show command line usage information.
	golf.Usage = func() {
		stderr("Use `%s --help` for more information.\n", ProgramName)
	}
}

var (
	optHelp    = golf.BoolP('h', "help", false, "Print command line help and exit")
	optQuiet   = golf.BoolP('q', "quiet", false, "Do not print intermediate errors to stderr")
	optVerbose = golf.BoolP('v', "verbose", false, "Print verbose output to stderr")

	optDelimiter = golf.StringP('d', "delimiter", "", "specify alternative field delimiter (empty string implies split on\n\twhitespace)")
	optField     = golf.StringP('f', "field", "", "Comma delimited list of field specifications to use as the histogram key.\n\tField numbering starts at 1. May include open ranges, such as '-3,5' for the\n\tfirst three fields, followed by the fifth field. The empty string implies\n\tentire line.")
	optFold      = golf.Bool("fold", false, "fold duplicate keys")
	optPercent   = golf.BoolP('p', "percentage", false, "show percentage")
	optRaw       = golf.Bool("raw", false, "Print keys and counts")
	optSortAsc   = golf.Bool("ascending", false, "print histogram in ascending order")
	optSortDesc  = golf.Bool("descending", false, "print histogram in descending order")
	optWidth     = golf.IntP('w', "width", 0, "width of output histogram. 0 implies use tty width")
)

func main() {
	golf.Parse()

	if *optHelp {
		// Show detailed help then exit, ignoring other possibly conflicting
		// options when '--help' is given.
		fmt.Printf(`histogram

Reads input from multiple files specified  on the command line or from standard
input when no files are specified.

SUMMARY:  histogram [options] [file1 [file2 ...]] [options]

USAGE: Not all options  may be used with all other  options. See below synopsis
for reference.

    histogram [--quiet | [--force | --verbose]]
              [--delimiter STRING] [--field INTEGER] [--fold]
              [--ascending | --descending]
              [--raw | [--percent | --width INTEGER]]
              [file1 [file2 ...]]

EXAMPLES:

    histogram < sample.txt
    histogram sample.txt
    last | histogram --field 1 --fold --descending

Command line options:
`)
		golf.PrintDefaults() // frustratingly, this only prints to stderr, and cannot change because it mimicks flag stdlib package
		return
	}

	if *optSortAsc && *optSortDesc {
		usage("cannot use both --ascending and --descending")
	}
	if *optRaw {
		if *optPercent {
			usage("cannot use both --raw and --percent")
		}
		if *optWidth > 0 {
			usage("cannot use both --raw and --width N")
		}
	} else if *optWidth == 0 {
		var err error
		*optWidth, _, err = gows.GetWinSize()
		if err != nil {
			warning("cannot get tty size (using raw output): %s", err)
			*optRaw = true
		}
	}

	fs, err := NewFieldSplitter(*optField, *optDelimiter)
	if err != nil {
		fatal(err)
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	sh := new(gohistogram.Strings)

	if err = ingest(ior, sh, fs); err != nil {
		fatal(err)
	}

	if *optFold {
		sh.FoldDuplicateKeys()
	}

	if *optSortDesc {
		sh.SortDescending()
	} else if *optSortAsc {
		sh.SortAscending()
	}

	if *optRaw {
		err = sh.PrintRaw()
	} else if *optPercent {
		err = sh.PrintWithPercent(*optWidth)
	} else {
		err = sh.Print(*optWidth)
	}
	if err != nil {
		fatal(err)
	}
}

func ingest(ior io.Reader, hist *gohistogram.Strings, fs *FieldSplitter) error {
	scanner := gobls.NewScanner(ior)
	for scanner.Scan() {
		// Remove line ending and split line into fields, then join into string
		key := fs.Select(strings.TrimRight(scanner.Text(), "\r\n"))

		// ignore empty string at the end of the input
		if len(key) > 0 {
			hist.Add(key)
		}
	}
	return scanner.Err()
}
