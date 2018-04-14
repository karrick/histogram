package main

import (
	"fmt"
	"testing"
)

func ExampleFieldSplitterFields() {
	f, err := NewFieldSplitter("2,4-5,8", "")
	if err != nil {
		panic(err) // for example use
	}
	fmt.Println(f.Fields("one two three four five six seven eight nine ten"))
	// Output: [two four five eight]
}

func ExampleFieldSplitterSelect() {
	f, err := NewFieldSplitter("2,4-5,8", "")
	if err != nil {
		panic(err) // for example use
	}
	fmt.Println(f.Select("one two three four five six seven eight nine ten"))
	// Output: two four five eight
}

func TestTextFieldsEmpty(t *testing.T) {
	tf, err := NewFieldSplitter("", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "one"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[2], "three"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsEmptyCommaDelimiter(t *testing.T) {
	tf, err := NewFieldSplitter("", ",")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "one two three"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsLeftAndRightMissing(t *testing.T) {
	_, err := NewFieldSplitter("-", "")
	if err == nil {
		t.Errorf("GOT: %v; WANT: %v", err, "non-nil")
	}
}

func TestTextFieldsInvertedRangeOrder(t *testing.T) {
	_, err := NewFieldSplitter("3-2", "")
	if err == nil {
		t.Errorf("GOT: %v; WANT: %v", err, "non-nil")
	}
}

func TestTextFieldsHyphen2(t *testing.T) {
	tf, err := NewFieldSplitter("-2", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "one"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFields2Hyphen(t *testing.T) {
	tf, err := NewFieldSplitter("2-", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 2; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "three"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFields2Comma4Hyphen5(t *testing.T) {
	tf, err := NewFieldSplitter("2,4-5", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three four five six")

	if got, want := len(fields), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "four"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[2], "five"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsIndexTooLarge(t *testing.T) {
	tf, err := NewFieldSplitter("4", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsLeftIndexExact(t *testing.T) {
	tf, err := NewFieldSplitter("2-", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two")

	if got, want := len(fields), 1; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsLeftIndexTooLarge(t *testing.T) {
	tf, err := NewFieldSplitter("4-", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 0; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsRightIndexTooLarge(t *testing.T) {
	tf, err := NewFieldSplitter("-4", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three")

	if got, want := len(fields), 3; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "one"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[2], "three"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}

func TestTextFieldsRightIndexExact(t *testing.T) {
	tf, err := NewFieldSplitter("-4", "")
	if err != nil {
		t.Errorf("GOT: %v; WANT: %v", err, nil)
	}

	fields := tf.Fields("one two three four")

	if got, want := len(fields), 4; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[0], "one"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[1], "two"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[2], "three"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := fields[3], "four"; got != want {
		t.Fatalf("GOT: %v; WANT: %v", got, want)
	}
}
