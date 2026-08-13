package main

import (
	"encoding/json"
	"fmt"
)

func printSliceInfo(label string, values []int) {
	fmt.Printf("%s: len=%d cap=%d\n", label, len(values), cap(values))
}

func filterEven(values []int) []int {
	result := make([]int, 0, len(values))
	for _, value := range values {
		if value%2 == 0 {
			result = append(result, value)
		}
	}
	return result
}

func printNames(names []string) {
	if len(names) == 0 {
		fmt.Println("Список имён пуст")
		return
	}
	fmt.Printf("Первое имя: %s\n", names[0])
	fmt.Printf("Последнее имя: %s\n", names[len(names)-1])
}

func compareNilAndEmpty() {
	var nilSlice []int
	emptySlice := []int{}

	fmt.Println("nil vs empty slice:")
	fmt.Printf("nilSlice == nil: %v\n", nilSlice == nil)
	fmt.Printf("emptySlice == nil: %v\n", emptySlice == nil)
	fmt.Printf("nilSlice len=%d cap=%d\n", len(nilSlice), cap(nilSlice))
	fmt.Printf("emptySlice len=%d cap=%d\n", len(emptySlice), cap(emptySlice))

	nilJSON, _ := json.Marshal(nilSlice)
	emptyJSON, _ := json.Marshal(emptySlice)
	fmt.Printf("nilSlice JSON: %s\n", string(nilJSON))
	fmt.Printf("emptySlice JSON: %s\n", string(emptyJSON))
}

func main() {
	literalSlice := []int{1, 2, 3, 4, 5}
	makeSlice := make([]int, 5, 8)
	arr := [5]int{10, 11, 12, 13, 14}
	arraySlice := arr[:]

	printSliceInfo("literalSlice", literalSlice)
	printSliceInfo("makeSlice", makeSlice)
	printSliceInfo("arraySlice", arraySlice)

	evenNumbers := filterEven([]int{1, 2, 3, 4, 5, 6, 7, 8})
	fmt.Printf("Чётные элементы: %v\n", evenNumbers)

	names := []string{}
	for _, name := range []string{"Анна", "Борис", "Вера", "Глеб", "Дина"} {
		names = append(names, name)
	}
	printNames(names)

	compareNilAndEmpty()
}
