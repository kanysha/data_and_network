package main

import "fmt"

func charCount(s string) map[rune]int {
	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}
	return counts
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func truncateWithEllipsis(s string, maxRunes int) string {
	runes := []rune(s)
	if maxRunes < 0 {
		maxRunes = 0
	}
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}

func main() {
	text := "привет"

	counts := charCount("hello")
	seen := make(map[rune]bool)
	ordered := make([]rune, 0, len(counts))
	for _, r := range []rune("hello") {
		if !seen[r] {
			seen[r] = true
			ordered = append(ordered, r)
		}
	}

	fmt.Printf("charCount(\"hello\") -> map[")
	for i, r := range ordered {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%c:%d", r, counts[r])
	}
	fmt.Println("]")
	fmt.Println("reverse('привет'):", reverse(text))
	fmt.Println("truncateWithEllipsis('привет мир', 5):", truncateWithEllipsis("привет мир", 5))
	fmt.Println("truncateWithEllipsis('коротко', 20):", truncateWithEllipsis("коротко", 20))
}
