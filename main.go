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
)

var (
	optDelimiter = golf.StringP('d', "delimiter", "", "specify alternative field delimiter (Default: empty string implies any whitespace")
	optField     = golf.IntP('f', "field", 0, "specify input field (Default: 0 implies entire line")
	optFold      = golf.Bool("fold", false, "fold duplicate keys")
	optHelp      = golf.BoolP('h', "help", false, "Print command line help and exit")
	optPercent   = golf.BoolP('p', "percentage", false, "show percentage")
	optSortAsc   = golf.Bool("ascending", false, "print histogram in ascending order")
	optSortDesc  = golf.Bool("descending", false, "print histogram in descending order")
	optWidth     = golf.IntP('w', "width", 80, "width of output histogram")
)

func main() {
	golf.Parse()

	if *optHelp {
		fmt.Fprintf(os.Stderr, "%s\n", filepath.Base(os.Args[0]))
		if *optHelp {
			fmt.Fprintf(os.Stderr, "\tgenerate and display a histogram from keys read from multiple lines\n\n")
			fmt.Fprintf(os.Stderr, "Reads input from multiple files specified on the command line or from\n")
			fmt.Fprintf(os.Stderr, "standard input when no files are specified.\n\n")
			golf.Usage()
		}
		exit(nil)
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	hist := new(gohistogram.Strings)
	err := ingest(ior, hist, *optField, *optDelimiter)
	if err != nil {
		exit(err)
	}

	if *optFold {
		hist.FoldDuplicateKeys()
	}

	if *optSortDesc {
		hist.SortDescending()
	} else if *optSortAsc {
		hist.SortAscending()
	}

	if *optPercent {
		exit(hist.PrintWithPercent(*optWidth))
	} else {
		exit(hist.Print(*optWidth))
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func ingest(ior io.Reader, hist *gohistogram.Strings, field int, delimiter string) error {
	scanner := gobls.NewScanner(ior)
	for scanner.Scan() {
		// split line into fields
		var fields []string
		key := scanner.Text()
		if field == 0 {
			key = strings.TrimSpace(key)
		} else {
			if delimiter == "" {
				fields = strings.Fields(key)
			} else {
				fields := strings.Split(key, delimiter)
				if len(fields) == 0 {
					continue
				}
			}
			if len(fields) <= field-1 {
				continue
			}
			key = fields[field-1]
		}
		// ignore empty string at the end of the input
		if len(key) > 0 {
			hist.Add(key)
		}
	}
	return scanner.Err()
}
