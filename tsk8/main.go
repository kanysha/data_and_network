package main

import (
	"fmt"
	"sort"
	"sync"
)

type Student struct {
	Name  string
	Grade int
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func unique(items []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(items))

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func groupStudentsByGrade(students []Student) map[int][]string {
	grouped := make(map[int][]string)

	for _, student := range students {
		grouped[student.Grade] = append(grouped[student.Grade], student.Name)
	}

	return grouped
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]string)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func printGroupedStudents(grouped map[int][]string) {
	grades := make([]int, 0, len(grouped))
	for grade := range grouped {
		grades = append(grades, grade)
	}
	sort.Ints(grades)

	for _, grade := range grades {
		fmt.Printf("%d:%v\n", grade, grouped[grade])
	}
}

func main() {
	fmt.Println("Exercise 1:")
	fmt.Println(unique([]string{"a", "b", "a", "c", "b"}))

	fmt.Println("\nExercise 2:")
	students := []Student{
		{"Alice", 5},
		{"Bob", 4},
		{"Carol", 5},
		{"Dave", 3},
		{"Eve", 4},
	}
	printGroupedStudents(groupStudentsByGrade(students))

	fmt.Println("\nExercise 3:")
	cache := NewCache()
	cache.Set("name", "Alice")
	cache.Set("city", "Moscow")
	value, ok := cache.Get("name")
	fmt.Println("Get name:", value, ok)
	fmt.Println("Keys:", cache.Keys())
	cache.Delete("city")
	fmt.Println("Keys after delete:", cache.Keys())
	cache.Clear()
	fmt.Println("Keys after clear:", cache.Keys())
}
