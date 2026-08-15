package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrInvalidInput = errors.New("invalid input")

func parseAge(s string) (int, error) {
	age, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse age %q: %w", s, err)
	}
	if age < 0 {
		return 0, fmt.Errorf("age %d is negative: %w", age, ErrInvalidInput)
	}
	return age, nil
}

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %q", e.Field)
}

func validateUser(name string, age int) error {
	if name == "" {
		return ErrInvalidInput
	}
	if age < 0 {
		return &ValidationError{Field: "age"}
	}
	return nil
}

func readAndCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan %q: %w", path, err)
	}
	return lineCount, nil
}

func demonstrateReadAndCount() error {
	file, err := os.CreateTemp("", "read-and-count-*.txt")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	path := file.Name()
	defer os.Remove(path)

	if _, err := file.WriteString("первая строка\nвторая строка\nтретья строка\n"); err != nil {
		file.Close()
		return fmt.Errorf("write temporary file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	count, err := readAndCount(path)
	if err != nil {
		return err
	}
	fmt.Println("Количество строк:", count)

	_, err = readAndCount(path + ".missing")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Файл не существует (проверено через errors.Is)")
		return nil
	}
	return fmt.Errorf("expected os.ErrNotExist, got %w", err)
}

func main() {
	fmt.Println("Задание 1:")
	for _, input := range []string{"25", "abc", "-5"} {
		age, err := parseAge(input)
		fmt.Printf("parseAge(%q) = %d, error = %v\n", input, age, err)
	}

	fmt.Println("\nЗадание 2:")
	if err := validateUser("", 20); errors.Is(err, ErrInvalidInput) {
		fmt.Println("Пустое имя:", err)
	}

	err := validateUser("Алиса", -1)
	if validationErr, ok := errors.AsType[*ValidationError](err); ok {
		fmt.Printf("Ошибка валидации поля %q\n", validationErr.Field)
	}

	fmt.Println("\nЗадание 3:")
	if err := demonstrateReadAndCount(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
