package gohistogram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type histItem struct {
	key   string
	count int
}

// Strings is a histogram of strings.
type Strings struct {
	countAfterFinalize int    // counts number of times Add called after last finalize
	keycount           int    // counts number of times Add called with previous key
	largestCount       int    // tracks the largest count of all the keys seen so far
	longestKey         int    // tracks the longest key length
	total              int    // counts the total number of items Add called
	key                string // tracks the key Add most recently called with
	items              []*histItem
}

// Add will add the specified key to the histogram.
func (hist *Strings) Add(key string) {
	hist.total++
	hist.countAfterFinalize++

	if hist.countAfterFinalize == 1 {
		// first key ever, or after finalize
		hist.key = key
		hist.keycount = 1
		return
	}

	if hist.key == key {
		hist.keycount++
		return
	}

	// new key

	hist.items = append(hist.items, &histItem{key: hist.key, count: hist.keycount})
	if hist.largestCount < hist.keycount {
		hist.largestCount = hist.keycount
	}
	if kl := len(hist.key); hist.longestKey < kl {
		hist.longestKey = kl
	}

	hist.key = key
	hist.keycount = 1
}

// finalize adds the current key and its count to the list of tracked items
func (hist *Strings) finalize() {
	if hist.countAfterFinalize > 0 {
		hist.items = append(hist.items, &histItem{key: hist.key, count: hist.keycount})
		if hist.largestCount < hist.keycount {
			hist.largestCount = hist.keycount
		}
		if kl := len(hist.key); hist.longestKey < kl {
			hist.longestKey = kl
		}
		hist.countAfterFinalize = 0
		hist.key = "" // not required, but release the previous key
	}
}

// Print displays the histogram with three columns: Key, Count, and a histogram of stars.
func (hist *Strings) Print(width int) error {
	const extra = 3 // space between key and count, space between count and histogram, plus 1 to keep from final column

	hist.finalize()

	if len(hist.items) > 0 {
		keyLength := hist.longestKey
		if l := len("Key"); keyLength < l {
			keyLength = l
		}
		countLength := len(strconv.FormatInt(int64(hist.largestCount), 10)) // width of the largest number
		if l := len("Count"); countLength < l {                             // ensure long enough for "Count"
			countLength = l
		}
		adjustedWidth := width - keyLength - countLength - extra
		if adjustedWidth < 1 {
			return fmt.Errorf("cannot print with fewer than %d columns", 1+width-adjustedWidth)
		}
		fmt.Printf("%-*s %*s (~%.3g per *)\n", keyLength, "Key", countLength, "Count", float64(hist.largestCount)/float64(adjustedWidth))
		for _, i := range hist.items {
			w := adjustedWidth * i.count / hist.largestCount
			fmt.Printf("%-*s %*d %s\n", keyLength, i.key, countLength, i.count, strings.Repeat("*", w))
		}
	}
	return nil
}

// PrintWithPercent displays the histogram with four columns: Key, Count, Percent, and a histogram of stars.
func (hist *Strings) PrintWithPercent(width int) error {
	const extra = 7 + 3 // len(percent), plus space between key and count, space between count and perc, space between perc and histogram, plus 1 to keep from final column

	hist.finalize()

	if len(hist.items) > 0 {
		keyLength := hist.longestKey
		if l := len("Key"); keyLength < l {
			keyLength = l
		}
		countLength := len(strconv.FormatInt(int64(hist.largestCount), 10)) // width of the largest number
		if l := len("Count"); countLength < l {                             // ensure long enough for "Count"
			countLength = l
		}
		adjustedWidth := width - keyLength - countLength - extra
		if adjustedWidth < 1 {
			return fmt.Errorf("cannot print with fewer than %d columns", 1+width-adjustedWidth)
		}
		perc := 100 / float64(hist.total)
		fmt.Printf("%-*s %*s Percent (~%.3g per *)\n", keyLength, "Key", countLength, "Count", float64(hist.largestCount)/float64(adjustedWidth))
		for _, i := range hist.items {
			w := adjustedWidth * i.count / hist.largestCount
			fmt.Printf("%-*s %*d % 7.2f %s\n", keyLength, i.key, countLength, i.count, (float64(i.count) * perc), strings.Repeat("*", w))
		}
	}
	return nil
}

// FoldDuplicateKeys aggregates counts of like keys in O(n) time.
func (hist *Strings) FoldDuplicateKeys() {
	hist.finalize()

	histItems := make([]*histItem, 0, len(hist.items))
	indexes := make(map[string]int)

	for _, item := range hist.items {
		histItemsIndex, ok := indexes[item.key]
		if ok {
			// already know about this key, but another reference to it
			c := histItems[histItemsIndex].count
			c += item.count
			histItems[histItemsIndex].count = c
			if hist.largestCount < c {
				hist.largestCount = c
			}
		} else {
			// have not seen this key before
			histItems = append(histItems, item)
			indexes[item.key] = len(histItems) - 1
		}
	}

	hist.items = histItems
}

func (hist *Strings) Len() int { return len(hist.items) }

func (hist *Strings) Less(i, j int) bool { return hist.items[i].count < hist.items[j].count }

func (hist *Strings) Swap(i, j int) {
	hist.items[j], hist.items[i] = hist.items[i], hist.items[j]
}

func (hist *Strings) SortAscending() {
	hist.finalize()
	sort.Sort(hist)
}

func (hist *Strings) SortDescending() {
	hist.finalize()
	sort.Sort(sort.Reverse(hist))
}
