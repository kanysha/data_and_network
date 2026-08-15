package main

import (
	"errors"
	"testing"
)

func TestFormatters(t *testing.T) {
	if got := (PlainText{Content: "Go"}).Format(); got != "Go" {
		t.Fatalf("PlainText.Format() = %q, want %q", got, "Go")
	}
	if got := (BoldText{Content: "Go"}).Format(); got != "**Go**" {
		t.Fatalf("BoldText.Format() = %q, want %q", got, "**Go**")
	}
}

func TestSumNumbers(t *testing.T) {
	got, err := sumNumbers(1, 2.9, "3")
	if err != nil {
		t.Fatalf("sumNumbers() returned error: %v", err)
	}
	if got != 6 {
		t.Fatalf("sumNumbers() = %d, want 6", got)
	}

	if _, err := sumNumbers(true); err == nil {
		t.Fatal("sumNumbers(bool) error = nil, want unsupported type error")
	}
}

func TestMemoryRepo(t *testing.T) {
	repo := NewMemoryRepo()
	if _, err := repo.FindByID(1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("FindByID() error = %v, want ErrNotFound", err)
	}
	if err := repo.Save(1, "value"); err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}
	if got, err := repo.FindByID(1); err != nil || got != "value" {
		t.Fatalf("FindByID() = %q, %v; want value, nil", got, err)
	}
}
