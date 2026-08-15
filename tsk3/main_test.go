package main

import (
	"reflect"
	"testing"
)

func TestAnalyzeLogs(t *testing.T) {
	got := analyzeLogs([]string{
		"GET /b 500",
		"GET /a 404",
		"POST /b 400",
		"invalid",
		"GET /ignored 200",
	})
	want := []pathErrorCount{{Path: "/b", Errors: 2}, {Path: "/a", Errors: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("analyzeLogs() = %v, want %v", got, want)
	}
}

func TestMergeSorted(t *testing.T) {
	got := mergeSorted([]int{1, 3, 5}, []int{2, 3, 4})
	want := []int{1, 2, 3, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mergeSorted() = %v, want %v", got, want)
	}
}

func TestMovingAvg(t *testing.T) {
	if got := movingAvg([]float64{1, 2, 3}, 5); got != nil {
		t.Fatalf("movingAvg(window > len) = %v, want nil", got)
	}

	got := movingAvg([]float64{1, 2, 3, 4}, 2)
	want := []float64{1.5, 2.5, 3.5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("movingAvg() = %v, want %v", got, want)
	}
}

func TestChunkBy(t *testing.T) {
	got := chunkBy([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("chunkBy() = %v, want %v", got, want)
	}
}
