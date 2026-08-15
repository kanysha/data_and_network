package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParseAge(t *testing.T) {
	if got, err := parseAge("25"); err != nil || got != 25 {
		t.Fatalf("parseAge(25) = %d, %v; want 25, nil", got, err)
	}
	if _, err := parseAge("-1"); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("parseAge(-1) error = %v, want ErrInvalidInput", err)
	}
	if _, err := parseAge("not-a-number"); err == nil {
		t.Fatal("parseAge(not-a-number) error = nil")
	}
}

func TestValidateUser(t *testing.T) {
	if err := validateUser("", 20); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("validateUser(empty name) error = %v, want ErrInvalidInput", err)
	}
	err := validateUser("Alice", -1)
	validationErr, ok := errors.AsType[*ValidationError](err)
	if !ok || validationErr.Field != "age" {
		t.Fatalf("validateUser(negative age) error = %v, want ValidationError for age", err)
	}
}

func TestReadAndCount(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lines.txt")
	if err := os.WriteFile(path, []byte("one\ntwo\nthree"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	got, err := readAndCount(path)
	if err != nil {
		t.Fatalf("readAndCount() returned error: %v", err)
	}
	if got != 3 {
		t.Fatalf("readAndCount() = %d, want 3", got)
	}
	if _, err := readAndCount(path + ".missing"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("readAndCount(missing) error = %v, want os.ErrNotExist", err)
	}
}
