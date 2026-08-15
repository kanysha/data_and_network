package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestUnique(t *testing.T) {
	got := unique([]string{"a", "b", "a", "c", "b"})
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unique() = %v, want %v", got, want)
	}
}

func TestGroupStudentsByGrade(t *testing.T) {
	students := []Student{{Name: "Alice", Grade: 5}, {Name: "Bob", Grade: 4}, {Name: "Carol", Grade: 5}}
	got := groupStudentsByGrade(students)
	want := map[int][]string{4: {"Bob"}, 5: {"Alice", "Carol"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupStudentsByGrade() = %v, want %v", got, want)
	}
}

func TestCache(t *testing.T) {
	cache := NewCache()
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			cache.Set(string(rune('a'+i)), "value")
		})
	}
	wg.Wait()

	if got := len(cache.Keys()); got != 20 {
		t.Fatalf("len(Keys()) = %d, want 20", got)
	}
	cache.Delete("a")
	if _, ok := cache.Get("a"); ok {
		t.Fatal("Get(a) found deleted key")
	}
	cache.Clear()
	if got := cache.Keys(); len(got) != 0 {
		t.Fatalf("Keys() after Clear = %v, want empty", got)
	}
}
