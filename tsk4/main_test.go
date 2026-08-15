package main

import (
	"reflect"
	"testing"
)

func TestFilterEven(t *testing.T) {
	got := filterEven([]int{-2, -1, 0, 1, 2})
	want := []int{-2, 0, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterEven() = %v, want %v", got, want)
	}
}
