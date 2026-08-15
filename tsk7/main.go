package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func wordFreq(text string) map[string]int {
	freq := make(map[string]int)
	for _, word := range strings.Fields(text) {
		freq[word]++
	}
	return freq
}

func reversePhonebook(book map[string]string) map[string]string {
	// При совпадении телефонов выбираем лексикографически последнее имя.
	// Явное правило делает результат независимым от случайного порядка map.
	names := slices.Sorted(maps.Keys(book))

	reverse := make(map[string]string)
	for _, name := range names {
		reverse[book[name]] = name
	}
	return reverse
}

func mergeMaps(m1, m2 map[string]int) map[string]int {
	result := make(map[string]int)
	for key, value := range m1 {
		result[key] = value
	}
	for key, value := range m2 {
		result[key] = value
	}
	return result
}

func main() {
	fmt.Println("wordFreq:", wordFreq("go go go is great"))

	phonebook := map[string]string{
		"Alice": "+7-900-111-22-33",
		"Bob":   "+7-900-444-55-66",
		"Carol": "+7-900-111-22-33",
	}
	fmt.Println("reversePhonebook:", reversePhonebook(phonebook))

	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 9, "c": 3}
	fmt.Println("mergeMaps:", mergeMaps(m1, m2))
}
