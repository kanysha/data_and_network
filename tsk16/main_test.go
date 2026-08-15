package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	got, err := countLines(strings.NewReader("one\ntwo\nthree"))
	if err != nil {
		t.Fatalf("countLines() returned error: %v", err)
	}
	if got != 3 {
		t.Fatalf("countLines() = %d, want 3", got)
	}
}

func TestHashReader(t *testing.T) {
	reader := NewHashReader(strings.NewReader("hello"))
	gotData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}
	if string(gotData) != "hello" {
		t.Fatalf("ReadAll() = %q, want hello", gotData)
	}
	wantHash := sha256.Sum256([]byte("hello"))
	if got, want := reader.Sum(), hex.EncodeToString(wantHash[:]); got != want {
		t.Fatalf("Sum() = %q, want %q", got, want)
	}
}

func TestCopyAndCount(t *testing.T) {
	var dst bytes.Buffer
	got, err := copyAndCount(&dst, strings.NewReader("test data"))
	if err != nil {
		t.Fatalf("copyAndCount() returned error: %v", err)
	}
	if got != 9 || dst.String() != "test data" {
		t.Fatalf("copyAndCount() = %d, %q; want 9, test data", got, dst.String())
	}
}

var errWrite = errors.New("write failed")

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
	return min(4, len(p)), errWrite
}

func TestCopyAndCountError(t *testing.T) {
	written, err := copyAndCount(failingWriter{}, strings.NewReader("test data"))
	if !errors.Is(err, errWrite) {
		t.Fatalf("copyAndCount() error = %v, want errWrite", err)
	}
	if written != 4 {
		t.Fatalf("copyAndCount() written = %d, want 4", written)
	}
}
