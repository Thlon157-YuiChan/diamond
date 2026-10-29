package model

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	cases := []struct {
		in   string
		want Time
		ok   bool
	}{
		{"00:00", 0, true},
		{"00:00:14", 1, true}, // snaps to 15s
		{"00:00:15", 1, true},
		{"01:00", 240, true},
		{"23:59:59", 5760, true}, // nearest 15s -> 翌日 00:00
		{"bad", 0, false},
		{"24:00", 0, false},
	}
	for _, c := range cases {
		got, ok := ParseTime(c.in)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("ParseTime(%q) = (%d,%v), want (%d,%v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

func TestTimeRoundTrip(t *testing.T) {
	tm := Time(100)
	if tm.String() == "" {
		t.Fatal("empty string")
	}
}
