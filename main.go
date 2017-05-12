package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	reverse   = flag.Bool("r", false, "reverse sort, so items are in ascending order")
	field     = flag.Int("f", 0, "specify input field (Default: 0 implies entire line")
	delimiter = flag.String("d", "", "specify alternative field delimiter (Default: empty string implies any whitespace")
)

func main() {
	flag.Parse()

	histogram := make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if *field == 0 {
			text = strings.TrimSpace(text)
		} else {
			if *delimiter == "" {
				text = strings.Fields(text)[*field-1]
			} else {
				fields := strings.Split(text, *delimiter)
				if len(fields) == 0 {
					continue
				}
				text = fields[*field-1]
			}
		}
		histogram[text] = histogram[text] + 1
	}
	if err := scanner.Err(); err != nil {
		bail(err)
	}

	// invert the hash
	items := make([]item, 0, len(histogram))
	for key, count := range histogram {
		i := sort.Search(len(items), func(i int) bool {
			if *reverse {
				return items[i].count > count
			} else {
				return items[i].count < count
			}
		})
		if i == len(items) {
			items = append(items, item{count: count, values: []string{key}})
			continue
		}
		if items[i].count == count {
			items[i].values = append(items[i].values, key)
			continue
		}
		f := item{count: count, values: []string{key}}
		items = append(items[:i], append([]item{f}, items[i:]...)...)
	}

	fmt.Println("Value Count")
	for _, foo := range items {
		for _, value := range foo.values {
			fmt.Println(foo.count, value)
		}
	}
}

func bail(err error) {
	fmt.Fprintf(os.Stderr, "%s", err)
	os.Exit(1)
}

type item struct {
	count  int
	values []string
}
