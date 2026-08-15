package tsk5

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func ParseCSV(line string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.FieldsPerRecord = -1
	record, err := r.Read()
	if errors.Is(err, io.EOF) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("parse CSV line: %w", err)
	}
	return record, nil
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return -1
	}, s)

	var result strings.Builder
	result.Grow(len(s))
	previousHyphen := false
	for _, r := range s {
		if r == '-' {
			if previousHyphen || result.Len() == 0 {
				continue
			}
			previousHyphen = true
		} else {
			previousHyphen = false
		}
		result.WriteRune(r)
	}

	return strings.TrimSuffix(result.String(), "-")
}

func ParseLog(line string) (level, timestamp, message string, err error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "[") {
		return "", "", "", ErrInvalidLogFormat
	}
	end := strings.Index(line, "]")
	if end == -1 {
		return "", "", "", ErrInvalidLogFormat
	}

	level = strings.TrimSpace(line[1:end])
	remainder := strings.TrimSpace(line[end+1:])
	parts := strings.SplitN(remainder, " - ", 2)
	if len(parts) != 2 {
		return "", "", "", ErrInvalidLogFormat
	}

	timestamp = strings.TrimSpace(parts[0])
	message = strings.TrimSpace(parts[1])
	if level == "" || timestamp == "" {
		return "", "", "", ErrInvalidLogFormat
	}
	return level, timestamp, message, nil
}

var ErrInvalidLogFormat = csvError("invalid log format")

type csvError string

func (e csvError) Error() string {
	return string(e)
}
