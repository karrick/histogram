package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/karrick/gobls"
)

type histItem struct {
	key   string
	count int
}

type histogram struct {
	total        int
	largestCount int
	count        int
	longestKey   int
	key          string
	items        []*histItem
}

func (hist *histogram) Ingest(ior io.Reader, field int, delimiter string) error {
	var atLeastOne bool
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
		if len(key) > 0 {
			hist.addKey(key)
			atLeastOne = true
			hist.total++
		}
	}
	if atLeastOne {
		hist.finalize()
	}
	return scanner.Err()
}

func (hist *histogram) addKey(key string) {
	if hist.key != key {
		if hist.key != "" {
			hist.items = append(hist.items, &histItem{key: hist.key, count: hist.count})
			if hist.largestCount < hist.count {
				hist.largestCount = hist.count
			}
			if kl := len(hist.key); hist.longestKey < kl {
				hist.longestKey = kl
			}
			hist.count = 1
		} else {
			hist.count = 1
		}
		hist.key = key
	} else {
		hist.count++
	}
}

func (hist *histogram) finalize() {
	hist.items = append(hist.items, &histItem{key: hist.key, count: hist.count})
	if hist.largestCount < hist.count {
		hist.largestCount = hist.count
	}
	if kl := len(hist.key); hist.longestKey < kl {
		hist.longestKey = kl
	}
}

func (hist *histogram) Print(width int) error {
	if len(hist.items) > 0 {
		countLength := len(strconv.FormatInt(int64(hist.largestCount), 10)) // width of the largest number
		if l := len("Count"); countLength < l {                             // ensure long enough for "Count"
			countLength = l
		}
		adjustedWidth := width - hist.longestKey - countLength - 3 // how many columns available for stars
		if adjustedWidth < 1 {
			return fmt.Errorf("cannot print with fewer than %d columns", width-adjustedWidth)
		}
		keyLength := hist.longestKey
		fmt.Printf("%*s %*s\n", keyLength, "Value", countLength, "Count")
		for _, i := range hist.items {
			w := adjustedWidth * i.count / hist.largestCount
			fmt.Printf("%*s %*d %s\n", keyLength, i.key, countLength, i.count, strings.Repeat("*", w))
		}
	}
	return nil
}

func (hist *histogram) PrintWithPercent(width int) error {
	if len(hist.items) > 0 {
		countLength := len(strconv.FormatInt(int64(hist.largestCount), 10)) // width of the largest number
		if l := len("Count"); countLength < l {                             // ensure long enough for "Count"
			countLength = l
		}
		adjustedWidth := width - hist.longestKey - countLength - 3 // how many columns available for stars
		if adjustedWidth < 1 {
			return fmt.Errorf("cannot print with fewer than %d columns", width-adjustedWidth)
		}
		keyLength := hist.longestKey
		fmt.Printf("%*s %*s Percent\n", keyLength, "Value", countLength, "Count")
		for _, i := range hist.items {
			w := adjustedWidth * i.count / hist.largestCount
			fmt.Printf("%*s %*d % 7.2f %s\n", keyLength, i.key, countLength, i.count, (100 * float64(i.count) / float64(hist.total)), strings.Repeat("*", w))
		}
	}
	return nil
}

// FoldDuplicateKeys aggregates counts of like keys in O(n) time.
func (hist *histogram) FoldDuplicateKeys() {
	histItems := make([]*histItem, 0, len(hist.items))
	indexes := make(map[string]int)

	for _, item := range hist.items {
		histItemsIndex, ok := indexes[item.key]
		if ok {
			// already know about this key, but another reference to it
			histItems[histItemsIndex].count += item.count
		} else {
			// have not seen this key before
			histItems = append(histItems, item)
			indexes[item.key] = len(histItems) - 1
		}
	}

	hist.items = histItems
}

func (hist *histogram) Len() int { return len(hist.items) }

func (hist *histogram) Less(i, j int) bool { return hist.items[i].count < hist.items[j].count }

func (hist *histogram) Swap(i, j int) {
	hist.items[j], hist.items[i] = hist.items[i], hist.items[j]
}

func (hist *histogram) SortAscending() {
	sort.Sort(hist)
}

func (hist *histogram) SortDescending() {
	sort.Sort(sort.Reverse(hist))
}
