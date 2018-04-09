package gohistogram

import (
	"testing"
)

func TestHistogramEmpty(t *testing.T) {
	hist := new(Strings)

	if got, want := len(hist.items), 0; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramSingle(t *testing.T) {
	hist := new(Strings)

	hist.Add("abc")
	hist.finalize()

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
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("def")
	hist.finalize()

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
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("abc")
	hist.Add("def")
	hist.finalize()

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
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("def")
	hist.Add("def")
	hist.finalize()

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
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("def")
	hist.Add("def")
	hist.Add("abc")
	hist.finalize()

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

func TestHistogramMultipleFinalize(t *testing.T) {
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("def")
	hist.finalize()
	hist.Add("def")
	hist.Add("abc")
	hist.finalize()

	if got, want := len(hist.items), 4; got != want {
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

	if got, want := hist.items[2].key, "def"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[2].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}

	if got, want := hist.items[3].key, "abc"; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
	if got, want := hist.items[3].count, 1; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}

func TestHistogramFold(t *testing.T) {
	hist := new(Strings)

	hist.Add("abc")
	hist.Add("def")
	hist.finalize()
	hist.Add("def")
	hist.Add("abc")
	hist.finalize()
	hist.FoldDuplicateKeys()

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
	if got, want := hist.items[1].count, 2; got != want {
		t.Errorf("GOT: %v; WANT: %v", got, want)
	}
}
