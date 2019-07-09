package main // import "github.com/karrick/histogram"

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karrick/gobls"
	"github.com/karrick/gohistogram"
	"github.com/karrick/golf"
	"github.com/karrick/gorill"
	"github.com/karrick/gows"
)

func init() {
	ProgramOneLineSummary = "generate and display histogram of keys"
	ProgramLongDescription = fmt.Sprintf("Reads input from multiple files specified on the command line or from standard\ninput when no files are specified.  For bug reports or feature requests, ask\nabout \"%s\" in the #oneliners channel.\n", ProgramName)
	ProgramVersion = "1.3.0"
	ProgramUsageExamples = fmt.Sprintf("SUMMARY:  %s [options] [file1 [file2 ...]]\n\nUSAGE:    Note that some options may not be used with other options.  See below\nsynopsis for reference.\n\n\t%s [--quiet | --verbose]\n\t\t[--delimiter STRING] [--field INTEGER] [--fold]\n\t\t[--ascending | --descending]\n\t\t[--raw | [--percent | --width INTEGER]]\n\nEXAMPLES:\n\tlast | %s --field 1 --fold --descending\n", ProgramName, ProgramName, ProgramName)
}

var (
	optDelimiter = golf.StringP('d', "delimiter", "", "specify alternative field delimiter (empty string implies split on\n\twhitespace)")
	optField     = golf.StringP('f', "field", "", "Comma delimited list of field specifications to use as the histogram key.\n\tField numbering starts at 1. May include open ranges, such as '-3,5' for the\n\tfirst three fields, followed by the fifth field. The empty string implies\n\tentire line.")
	optFold      = golf.Bool("fold", false, "fold duplicate keys")
	optPercent   = golf.BoolP('p', "percentage", false, "show percentage")
	optRaw       = golf.Bool("raw", false, "Print keys and counts")
	optSortAsc   = golf.Bool("ascending", false, "print histogram in ascending order")
	optSortDesc  = golf.Bool("descending", false, "print histogram in descending order")
	optWidth     = golf.IntP('w', "width", 0, "width of output histogram. 0 implies use tty width")
)

func cmd() error {
	var err error

	if *optSortAsc && *optSortDesc {
		return ErrUsage("cannot use both --ascending and --descending")
	}
	if *optRaw {
		if *optPercent {
			return ErrUsage("cannot use both --raw and --percent")
		}
		if *optWidth > 0 {
			return ErrUsage("cannot use both --raw and --width N")
		}
	} else if *optWidth == 0 {
		*optWidth, _, err = gows.GetWinSize()
		if err != nil {
			warning("cannot get tty size (using raw output): %s", err)
			*optRaw = true
		}
	}

	fs, err := NewFieldSplitter(*optField, *optDelimiter)
	if err != nil {
		return err
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	sh := new(gohistogram.Strings)

	if err = ingest(ior, sh, fs); err != nil {
		return err
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
		return sh.PrintRaw()
	} else if *optPercent {
		return sh.PrintWithPercent(*optWidth)
	}
	return sh.Print(*optWidth)
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
