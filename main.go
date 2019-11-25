package main // import "github.com/karrick/histogram"

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karrick/gobls"
	"github.com/karrick/gohistogram"
	"github.com/karrick/golf"
	"github.com/karrick/gologs"
	"github.com/karrick/gorill"
	"github.com/karrick/gows"
)

var log *gologs.Logger

func init() {
	var err error
	log, err = gologs.New(os.Stderr, "{program}: {message}")
	if err != nil {
		panic(err)
	}
	golf.Usage = func() {
		log.User("Use `--help` for more information.\n")
	}
}

func fatal(f string, args ...interface{}) {
	log.User(f, args...)
	os.Exit(1)
}

func usage(f string, args ...interface{}) {
	log.User(f, args...)
	golf.Usage()
	os.Exit(2)
}

var (
	optHelp    = golf.BoolP('h', "help", false, "Print command line help and exit")
	optDebug   = golf.Bool("debug", false, "Print debug output to stderr")
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

    histogram [--debug | --verbose]
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
		golf.PrintDefaultsTo(os.Stdout) // frustratingly, this only prints to stderr, and cannot change because it mimicks flag stdlib package
		return
	}

	if *optDebug {
		log.SetDev()
	} else if *optVerbose {
		log.SetAdmin()
	} else {
		log.SetUser()
	}

	if *optSortAsc && *optSortDesc {
		usage("cannot use both --ascending and --descending")
	}

	fs, err := NewFieldSplitter(*optField, *optDelimiter)
	if err != nil {
		usage("%s", err)
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
			log.Admin("cannot get tty size (using raw output): %s", err)
			*optRaw = true
		}
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	sh := new(gohistogram.Strings)

	scanner := gobls.NewScanner(ior)
	for scanner.Scan() {
		// Remove line ending and split line into fields, then join into string
		key := fs.Select(strings.TrimRight(scanner.Text(), "\r\n"))
		// ignore empty string at the end of the input
		if len(key) > 0 {
			sh.Add(key)
		}
	}
	if err := scanner.Err(); err != nil {
		fatal("%s\n", err)
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
		fatal("%s\n", err)
	}
}
