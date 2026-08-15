package main

import "testing"

const benchmarkSize = 10_000

func TestFirstTenMemoryBehavior(t *testing.T) {
	source := make([]byte, 20)
	shared := firstTenShared(source)
	copied := firstTenCopied(source)

	shared[0] = 1
	if source[0] != 1 {
		t.Fatal("firstTenShared() should share the source array")
	}

	copied[1] = 2
	if source[1] != 0 {
		t.Fatal("firstTenCopied() should not share the source array")
	}
}

func BenchmarkAppendWithoutMake(b *testing.B) {
	for b.Loop() {
		var values []int
		for value := 0; value < benchmarkSize; value++ {
			values = append(values, value)
		}
	}
}

func BenchmarkAppendWithMake(b *testing.B) {
	for b.Loop() {
		values := make([]int, 0, benchmarkSize)
		for value := 0; value < benchmarkSize; value++ {
			values = append(values, value)
		}
	}
}
