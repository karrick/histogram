package main

import (
	"fmt"
	"io"
	"os"

	"github.com/karrick/golf"
	"github.com/karrick/gorill"
)

var (
	optDelimiter  = golf.StringP('d', "delimiter", "", "specify alternative field delimiter (Default: empty string implies any whitespace")
	optField      = golf.IntP('f', "field", 0, "specify input field (Default: 0 implies entire line")
	optFold       = golf.Bool("fold", false, "fold duplicate keys")
	optPercentage = golf.BoolP('p', "percentage", false, "show percentage")
	optSortAsc    = golf.Bool("ascending", false, "print histogram in ascending order")
	optSortDesc   = golf.Bool("descending", false, "print histogram in descending order")
	optWidth      = golf.IntP('w', "width", 80, "width of output histogram")
)

func main() {
	golf.Parse()

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	hist := new(histogram)
	err := hist.Ingest(ior, *optField, *optDelimiter)
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

	if *optPercentage {
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
