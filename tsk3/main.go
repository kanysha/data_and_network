package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type pathErrorCount struct {
	Path   string
	Errors int
}

func analyzeLogs(logLines []string) []pathErrorCount {
	counts := make(map[string]int)

	for _, line := range logLines {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		status, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			continue
		}

		path := parts[len(parts)-2]
		if status >= 400 {
			counts[path]++
		}
	}

	items := make([]pathErrorCount, 0, len(counts))
	for path, count := range counts {
		items = append(items, pathErrorCount{Path: path, Errors: count})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Errors == items[j].Errors {
			return items[i].Path < items[j].Path
		}
		return items[i].Errors > items[j].Errors
	})

	return items
}

func mergeSorted(a, b []int) []int {
	result := make([]int, 0, len(a)+len(b))
	left, right := 0, 0

	for left < len(a) && right < len(b) {
		if a[left] <= b[right] {
			result = append(result, a[left])
			left++
		} else {
			result = append(result, b[right])
			right++
		}
	}

	for left < len(a) {
		result = append(result, a[left])
		left++
	}

	for right < len(b) {
		result = append(result, b[right])
		right++
	}

	return result
}

func movingAvg(data []float64, window int) []float64 {
	if window <= 0 || window > len(data) {
		return nil
	}

	result := make([]float64, 0, len(data)-window+1)
	sum := 0.0

	for i, value := range data {
		sum += value
		if i >= window {
			sum -= data[i-window]
		}

		if i >= window-1 {
			result = append(result, sum/float64(window))
		}
	}

	return result
}

func chunkBy(items []int, chunkSize int) [][]int {
	if chunkSize <= 0 {
		return nil
	}

	chunks := make([][]int, 0, (len(items)+chunkSize-1)/chunkSize)
	for start := 0; start < len(items); start += chunkSize {
		end := start + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[start:end])
	}
	return chunks
}

func main() {
	fmt.Println("Задача 1: анализ логов")
	logLines := []string{
		"GET /home 200",
		"GET /home 404",
		"GET /about 500",
		"POST /home 404",
		"GET /about 200",
		"GET /contact 503",
		"POST /contact 400",
		"GET /contact 200",
	}

	for _, item := range analyzeLogs(logLines) {
		fmt.Printf("%s -> %d\n", item.Path, item.Errors)
	}

	fmt.Println("\nЗадача 2: слияние отсортированных массивов")
	a := []int{1, 3, 5, 7, 9}
	b := []int{2, 4, 6, 8, 10}
	fmt.Println(mergeSorted(a, b))

	fmt.Println("\nЗадача 3: скользящее окно")
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(movingAvg(data, 3))

	fmt.Println("\nБонус: разбиение на части")
	fmt.Println(chunkBy([]int{1, 2, 3, 4, 5, 6, 7}, 3))
}
