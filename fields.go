package main

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldSplitter parses strings into zero, one, or more fields according to the
// field specification string and a field delimiter. Note that field numbers
// start at 1. Also note that the number of fields returned may not match the
// number of field specifiers where there are ranges, or when a particular line
// has fewer fields than a particular specifier specifies. Ranges may be open
// ended or close ended.
//
//     func ExampleFieldSplitter() {
//         f, err := NewFieldSplitter("2,4-5,8", "")
//         if err != nil {
//             panic(err) // for example use
//         }
//         fmt.Println(f.Fields("one two three four five six seven eight nine ten"))
//         // Output: [two four five eight]
//     }
type FieldSplitter struct {
	fieldDelimiter     string                    // used to split string. when empty, splits on whitespace
	fps                []func([]string) []string // field picker functions
	fieldCountEstimate int                       // estimate number of fields each Fields() method will return
}

// NewFieldSplitter returns a FieldSplitter.
func NewFieldSplitter(commaDelimitedSpecs, fieldDelimiter string) (*FieldSplitter, error) {
	fs := &FieldSplitter{fieldDelimiter: fieldDelimiter, fieldCountEstimate: 1}

	if commaDelimitedSpecs == "" {
		return fs, nil
	}

	var err error

	specs := strings.Split(commaDelimitedSpecs, ",")

	fs.fps = make([]func([]string) []string, len(specs)) // we know exactly how many functions to call

	for i, spec := range specs {
		// Because no such thing as a negative field number, when first byte is
		// hyphen, skip Atoi attempt.
		if spec[0] != '-' {
			if someInt, err := strconv.Atoi(spec); err == nil {
				// This field spec was a simple integer, e.g., "3".
				fs.fps[i] = func(ss []string) []string {
					if len(ss) < someInt {
						return nil
					}
					return []string{ss[someInt-1]}
				}
				fs.fieldCountEstimate++
				continue // next field specification
			}
		}
		// When field specification is not a simple number, then split on hyphen
		// and expect 2 components.
		components := strings.Split(spec, "-")
		if len(components) != 2 {
			return nil, fmt.Errorf("cannot parse field specification: %q", spec)
		}
		// When not token values, left and right specify the starting field and
		// ending field, inclusive.
		left := -1
		right := -1
		// Attempt to parse the left side of the field specification.
		if components[0] != "" {
			left, err = strconv.Atoi(components[0])
			if err != nil {
				return nil, fmt.Errorf("cannot parse left side of field specification: %q", components[0])
			}
			if left < 1 {
				return nil, fmt.Errorf("cannot use zero or negative field specification: %q", components[0])
			}
		}
		// Attempt to parse the right side of the field specification.
		if components[1] != "" {
			right, err = strconv.Atoi(components[1])
			if err != nil {
				return nil, fmt.Errorf("cannot parse right side of field specification: %q", components[1])
			}
			if right < 1 {
				return nil, fmt.Errorf("cannot use zero or negative field specification: %q", components[1])
			}
		}
		// Only left was missing, only right was missing, both were missing
		// (error), or left > right (error), or left <= right.
		if right == -1 {
			if left == -1 {
				return nil, fmt.Errorf("cannot use invalid field specification: %q", spec)
			}
			// L-
			fs.fps[i] = func(ss []string) []string {
				if len(ss) < left {
					return nil
				}
				return ss[left-1:]
			}
			fs.fieldCountEstimate += left
		} else if left == -1 {
			// -R
			fs.fps[i] = func(ss []string) []string {
				if len(ss) < right {
					return ss
				}
				return ss[:right]
			}
			// Expect at least one field, which is not entirely accurate,
			// because field specification could be "5-", and there might be 10
			// fields in a particular input string.
			fs.fieldCountEstimate++
		} else if left > right {
			return nil, fmt.Errorf("left side cannot be less than right side of field specification: %q", specs)
		} else {
			// L-R
			fs.fps[i] = func(ss []string) []string {
				return ss[left-1 : right]
			}
			fs.fieldCountEstimate += 1 + right - left
		}
	}
	return fs, nil
}

// Fields splits the input string into a slice of strings based on the
// configured delimiter and the configured field specification string. See
// examples for this data type.
func (fs *FieldSplitter) Fields(s string) []string {
	var fields []string // fields

	if fs.fieldDelimiter != "" {
		fields = strings.Split(s, fs.fieldDelimiter)
	} else {
		fields = strings.Fields(s)
	}

	if len(fs.fps) == 0 {
		return fields // when no field specifiers, return slice of all fields
	}

	// Presize will not always be accurate, e.g., when a field spec is "5-", and
	// there are more than 5 fields in a particular string, however this handles
	// most cases without a second memory allocation.
	rs := make([]string, 0, fs.fieldCountEstimate)

	for _, fp := range fs.fps {
		// Recall that field picker function might return 0, 1, or more fields.
		rs = append(rs, fp(fields)...)
	}

	return rs
}

// Select returns a string representing only the selected fields from the input
// string. It is equivalent to splitting the input string on the delimiter,
// collecting the fields specified by the field specifications, then joining the
// resultant fields again with the field delimiter.
func (fs *FieldSplitter) Select(s string) string {
	if fs.fieldDelimiter != "" {
		return strings.Join(fs.Fields(s), fs.fieldDelimiter)
	}
	return strings.Join(fs.Fields(s), " ")
}
