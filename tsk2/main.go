package main

import (
	"fmt"
	"strings"
)

func removeAll(s []int, target int) []int {
	write := 0
	for _, v := range s {
		if v != target {
			s[write] = v
			write++
		}
	}
	return s[:write]
}

func insertSorted(s []int, v int) []int {
	idx := len(s)
	for i, existing := range s {
		if existing > v {
			idx = i
			break
		}
	}

	result := append(s, 0)
	copy(result[idx+1:], result[idx:])
	result[idx] = v
	return result
}

func matrixTranspose(matrix [][]int) [][]int {
	if len(matrix) == 0 {
		return [][]int{}
	}

	rows := len(matrix)
	cols := len(matrix[0])
	transposed := make([][]int, cols)
	for i := range transposed {
		transposed[i] = make([]int, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}

	return transposed
}

func csvLine(fields []string) string {
	var b strings.Builder
	for i, field := range fields {
		if i > 0 {
			b.WriteByte(',')
		}
		if needsQuotes(field) {
			b.WriteByte('"')
			for _, r := range field {
				if r == '"' {
					b.WriteString("\"\"")
				} else {
					b.WriteRune(r)
				}
			}
			b.WriteByte('"')
		} else {
			b.WriteString(field)
		}
	}
	return b.String()
}

func needsQuotes(field string) bool {
	for _, r := range field {
		if r == ',' || r == '"' || r == '\n' || r == '\r' {
			return true
		}
	}
	return false
}

func main() {
	s := []int{1, 2, 3, 2, 4, 2, 5}
	target := 2
	result := removeAll(s, target)
	fmt.Println("После удаления", target, ":", result)

	// Example usage of insertSorted
	sorted := []int{1, 3, 5, 7}
	valueToInsert := 4
	newSorted := insertSorted(sorted, valueToInsert)
	fmt.Println("После вставки", valueToInsert, ":", newSorted)

	// Example usage of matrixTranspose
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	transposedMatrix := matrixTranspose(matrix)
	fmt.Println("Транспонированная матрица:")
	for _, row := range transposedMatrix {
		fmt.Println(row)
	}

	// Example usage of csvLine
	fields := []string{"field1", "field,2", "field\"3", "field\n4"}
	csv := csvLine(fields)
	fmt.Println("CSV линия:", csv)
}
