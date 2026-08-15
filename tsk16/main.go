package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
)

func countLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("count lines: %w", err)
	}
	return count, nil
}

type HashReader struct {
	source io.Reader
	hash   hash.Hash
}

func NewHashReader(r io.Reader) *HashReader {
	return &HashReader{
		source: r,
		hash:   sha256.New(),
	}
}

func (h *HashReader) Read(p []byte) (int, error) {
	n, err := h.source.Read(p)
	if n > 0 {
		if _, hashErr := h.hash.Write(p[:n]); hashErr != nil {
			return n, fmt.Errorf("update hash: %w", hashErr)
		}
	}
	return n, err
}

func (h *HashReader) Sum() string {
	return hex.EncodeToString(h.hash.Sum(nil))
}

func copyAndCount(dst io.Writer, src io.Reader) (int64, error) {
	written, err := io.Copy(dst, src)
	if err != nil {
		return written, fmt.Errorf("copy data: %w", err)
	}
	return written, nil
}

func main() {
	fmt.Println("Задание 1:")
	lines, err := countLines(strings.NewReader("line1\nline2\nline3\n"))
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Строк:", lines)
	}

	fmt.Println("\nЗадание 2:")
	hashReader := NewHashReader(strings.NewReader("hello world"))
	data, err := io.ReadAll(hashReader)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Данные:", string(data))
		fmt.Println("SHA-256:", hashReader.Sum())
	}

	fmt.Println("\nЗадание 3:")
	var destination bytes.Buffer
	copied, err := copyAndCount(&destination, strings.NewReader("test data"))
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Скопировано: %d байт, содержимое: %s\n", copied, destination.String())
	}
}
