package main

import "testing"

const benchmarkSize = 10_000

func BenchmarkAppendWithoutMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var values []int
		for value := 0; value < benchmarkSize; value++ {
			values = append(values, value)
		}
	}
}

func BenchmarkAppendWithMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		values := make([]int, 0, benchmarkSize)
		for value := 0; value < benchmarkSize; value++ {
			values = append(values, value)
		}
	}
}
