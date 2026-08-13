package project5

import (
	"encoding/csv"
	"io"
	"strings"
	"unicode"
)

func ParseCSV(line string) []string {
	r := csv.NewReader(strings.NewReader(line))
	r.FieldsPerRecord = -1
	record, err := r.Read()
	if err == io.EOF {
		return []string{}
	}
	if err != nil {
		return []string{}
	}
	return record
}

func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return -1
	}, s)

	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	s = strings.Trim(s, "-")
	return s
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
