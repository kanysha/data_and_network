package main

import (
	"fmt"
)

func firstTenShared(s []byte) []byte {
	a := s[:10]
	return a
}

func firstTenCopied(s []byte) []byte {
	a := make([]byte, 10)
	copy(a, s[:10])
	return a
}

func TaskOne() {
	s := make([]int, 3, 10)
	fmt.Printf("Оригинальная длина: %d\nОригинальная емкость: %d\n", len(s), cap(s))

	a := s[0:2]
	fmt.Printf("Длина подмассива: %d\nЕмкость подмассива: %d\n", len(a), cap(a))

	b := s[1:5]
	fmt.Printf("Длина подмассива: %d\nЕмкость подмассива: %d\n", len(b), cap(b))

	c := s[:10]
	fmt.Printf("Длина подмассива: %d\nЕмкость подмассива: %d\n", len(c), cap(c))

	fmt.Println("Для среза s[start:end] длина равна end-start, а ёмкость — cap(s)-start.")
}

func TaskTwo() {
	s := make([]byte, 1<<20)
	copy(s[0:10], []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	a := firstTenShared(s)
	fmt.Printf("Слайс первых 10 байт:%v\n", a)
	b := firstTenCopied(s)
	fmt.Printf("Слайс первых 10 байт (с исправленной утечкой):%v\n", b)
}

func main() {
	TaskOne()
	TaskTwo()
	fmt.Println("Бенчмарки: go test ./tsk1 -bench=Append -benchmem")
}
