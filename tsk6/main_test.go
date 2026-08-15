package main

import (
	"reflect"
	"testing"
)

func TestCharCount(t *testing.T) {
	got := charCount("го-го")
	want := map[rune]int{'г': 2, 'о': 2, '-': 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("charCount() = %v, want %v", got, want)
	}
}

func TestReverse(t *testing.T) {
	if got, want := reverse("привет🙂"), "🙂тевирп"; got != want {
		t.Fatalf("reverse() = %q, want %q", got, want)
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	tests := []struct {
		name string
		text string
		max  int
		want string
	}{
		{name: "unicode", text: "привет", max: 3, want: "при..."},
		{name: "unchanged", text: "go", max: 2, want: "go"},
		{name: "negative", text: "go", max: -1, want: "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateWithEllipsis(tt.text, tt.max); got != tt.want {
				t.Fatalf("truncateWithEllipsis() = %q, want %q", got, tt.want)
			}
		})
	}
}
