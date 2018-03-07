package main

import (
	"bytes"
	"testing"
)

func TestHistogramEmpty(t *testing.T) {
	hist := new(histogram)
	bb := new(bytes.Buffer)

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 0; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramSingle(t *testing.T) {
	hist := new(histogram)
	bb := bytes.NewReader([]byte("abc\n"))

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	item0 := hist.items[0]

	if got, want := item0.key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := item0.count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramTwoSingles(t *testing.T) {
	hist := new(histogram)
	bb := bytes.NewReader([]byte("abc\ndef\n"))

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[0].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[0].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[1].key, "def"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[1].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramTwoFirstOneSecond(t *testing.T) {
	hist := new(histogram)
	bb := bytes.NewReader([]byte("abc\nabc\n\ndef\n"))

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[0].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[0].count, 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[1].key, "def"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[1].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramOneFirstTwoSecond(t *testing.T) {
	hist := new(histogram)
	bb := bytes.NewReader([]byte("abc\ndef\n\ndef\n"))

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[0].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[0].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[1].key, "def"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[1].count, 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramOneFirstTwoSecondRepeatFirst(t *testing.T) {
	hist := new(histogram)
	bb := bytes.NewReader([]byte("abc\ndef\n\ndef\nabc\n"))

	err := hist.Ingest(bb, 0, " ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(hist.items), 3; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[0].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[0].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[1].key, "def"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[1].count, 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[2].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[2].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

}
