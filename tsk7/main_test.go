package main

import (
	"reflect"
	"testing"
)

func TestWordFreq(t *testing.T) {
	got := wordFreq("go go\tgo\nis  great")
	want := map[string]int{"go": 3, "is": 1, "great": 1}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("wordFreq() = %v, want %v", got, want)
	}
}

func TestReversePhonebook(t *testing.T) {
	phonebook := map[string]string{
		"Alice": "+7-900-111-22-33",
		"Bob":   "+7-900-444-55-66",
		"Carol": "+7-900-111-22-33",
	}

	got := reversePhonebook(phonebook)
	want := map[string]string{
		"+7-900-111-22-33": "Carol",
		"+7-900-444-55-66": "Bob",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("reversePhonebook() = %v, want %v", got, want)
	}
}

func TestMergeMaps(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 9, "c": 3}

	got := mergeMaps(m1, m2)
	want := map[string]int{"a": 1, "b": 9, "c": 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mergeMaps() = %v, want %v", got, want)
	}
}
