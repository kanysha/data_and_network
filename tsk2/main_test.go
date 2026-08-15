package main

import (
	"reflect"
	"testing"
)

func TestRemoveAll(t *testing.T) {
	got := removeAll([]int{1, 2, 3, 2, 4}, 2)
	want := []int{1, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("removeAll() = %v, want %v", got, want)
	}
}

func TestInsertSorted(t *testing.T) {
	got := insertSorted([]int{1, 3, 5}, 3)
	want := []int{1, 3, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("insertSorted() = %v, want %v", got, want)
	}
}

func TestMatrixTranspose(t *testing.T) {
	got := matrixTranspose([][]int{{1, 2, 3}, {4, 5, 6}})
	want := [][]int{{1, 4}, {2, 5}, {3, 6}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("matrixTranspose() = %v, want %v", got, want)
	}
}

func TestCSVLine(t *testing.T) {
	got := csvLine([]string{"plain", "with,comma", `with"quote`, "with\nnewline"})
	want := "plain,\"with,comma\",\"with\"\"quote\",\"with\nnewline\""
	if got != want {
		t.Fatalf("csvLine() = %q, want %q", got, want)
	}
}
