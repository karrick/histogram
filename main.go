package main

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

var (
	optDelimiter = golf.StringP('d', "delimiter", "", "specify alternative field delimiter (empty string implies split on\n\twhitespace)")
	optField     = golf.StringP('f', "field", "", "Comma delimited list of field specifications to use as the histogram key.\n\tField numbering starts at 1. May include open ranges, such as '-3,5' for the\n\tfirst three fields, followed by the fifth field. The empty string implies\n\tentire line.")
	optFold      = golf.Bool("fold", false, "fold duplicate keys")
	optHelp      = golf.BoolP('h', "help", false, "Print command line help and exit")
	optPercent   = golf.BoolP('p', "percentage", false, "show percentage")
	optSortAsc   = golf.Bool("ascending", false, "print histogram in ascending order")
	optSortDesc  = golf.Bool("descending", false, "print histogram in descending order")
	optWidth     = golf.IntP('w', "width", 0, "width of output histogram. 0 implies use tty width")
)

func main() {
	golf.Parse()

	if *optHelp {
		fmt.Fprintf(os.Stderr, "%s\n", filepath.Base(os.Args[0]))
		if *optHelp {
			fmt.Fprintf(os.Stderr, "\tGenerate and display a histogram of keys.\n\n")
			fmt.Fprintf(os.Stderr, "Reads input from multiple files specified on the command line or from\n")
			fmt.Fprintf(os.Stderr, "standard input when no files are specified.\n\n")
			golf.Usage()
		}
		exit(nil)
	}

	fs, err := NewFieldSplitter(*optField, *optDelimiter)
	if err != nil {
		exit(err)
	}

	if *optWidth == 0 {
		// ignore error; if cannot get window size, use 80
		*optWidth, _, _ = gows.GetWinSize()
		if *optWidth == 0 {
			*optWidth = 80
		}
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	sh := new(gohistogram.Strings)

	if err = ingest(ior, sh, fs); err != nil {
		exit(err)
	}

	if *optFold {
		sh.FoldDuplicateKeys()
	}

	if *optSortDesc {
		sh.SortDescending()
	} else if *optSortAsc {
		sh.SortAscending()
	}

	if *optPercent {
		exit(sh.PrintWithPercent(*optWidth))
	} else {
		exit(sh.Print(*optWidth))
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
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
