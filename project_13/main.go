package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Formatter describes a value that can render itself as text.
type Formatter interface {
	Format() string
}

type PlainText struct {
	Content string
}

func (p PlainText) Format() string {
	return p.Content
}

type BoldText struct {
	Content string
}

func (b BoldText) Format() string {
	return "**" + b.Content + "**"
}

func printFormatted(f Formatter) {
	fmt.Println(f.Format())
}

// sumNumbers sums supported values. float64 values are converted to int,
// as required by the result type in the exercise.
func sumNumbers(values ...any) (int, error) {
	total := 0
	for _, value := range values {
		switch v := value.(type) {
		case int:
			total += v
		case float64:
			total += int(v)
		case string:
			n, err := strconv.Atoi(v)
			if err != nil {
				return 0, fmt.Errorf("convert %q to int: %w", v, err)
			}
			total += n
		default:
			return 0, fmt.Errorf("unsupported number type %T", value)
		}
	}
	return total, nil
}

type Repository interface {
	FindByID(id int) (string, error)
	Save(id int, data string) error
}

var ErrNotFound = errors.New("not found")

type MemoryRepo struct {
	data map[int]string
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{data: make(map[int]string)}
}

func (m *MemoryRepo) FindByID(id int) (string, error) {
	data, ok := m.data[id]
	if !ok {
		return "", fmt.Errorf("find record %d: %w", id, ErrNotFound)
	}
	return data, nil
}

func (m *MemoryRepo) Save(id int, data string) error {
	if m.data == nil {
		m.data = make(map[int]string)
	}
	m.data[id] = data
	return nil
}

func processRepo(repo Repository) error {
	if err := repo.Save(1, "Hello"); err != nil {
		return fmt.Errorf("save record: %w", err)
	}

	value, err := repo.FindByID(1)
	if err != nil {
		return fmt.Errorf("find saved record: %w", err)
	}
	fmt.Println("Найдено:", value)
	return nil
}

func main() {
	fmt.Println("Задание 1:")
	printFormatted(PlainText{Content: "Привет"})
	printFormatted(BoldText{Content: "Привет"})

	fmt.Println("\nЗадание 2:")
	total, err := sumNumbers(1, 2.0, "3")
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Сумма:", total)
	}

	fmt.Println("\nЗадание 3:")
	if err := processRepo(NewMemoryRepo()); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
