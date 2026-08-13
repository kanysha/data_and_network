package project5

import (
	"reflect"
	"testing"
)

func TestParseCSV(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []string
	}{
		{"simple", "a,b,c", []string{"a", "b", "c"}},
		{"quoted comma", "1,\"two, too\",3", []string{"1", "two, too", "3"}},
		{"quoted field", "\"a,b\",c", []string{"a,b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCSV(tt.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseCSV(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "Hello World!", "hello-world"},
		{"special chars", "Go & Rust", "go-rust"},
		{"spaces", "  A  B  C  ", "a-b-c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Fatalf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseLog(t *testing.T) {
	line := "[INFO] 2026-08-04T12:00:00Z - Started"
	level, timestamp, message, err := ParseLog(line)
	if err != nil {
		t.Fatalf("ParseLog returned error: %v", err)
	}
	if level != "INFO" || timestamp != "2026-08-04T12:00:00Z" || message != "Started" {
		t.Fatalf("parseLog(%q) = %q, %q, %q", line, level, timestamp, message)
	}
}

func TestParseLogInvalid(t *testing.T) {
	_, _, _, err := ParseLog("INFO 2026-08-04T12:00:00Z - Started")
	if err == nil {
		t.Fatal("expected error for invalid log format")
	}
}
