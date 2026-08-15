package main

import (
	"fmt"

	"github.com/kanysha/data_and_network/tsk5"
)

func main() {
	csvLine := "1,\"two, too\",3"
	parsed, err := tsk5.ParseCSV(csvLine)
	if err != nil {
		fmt.Printf("parseCSV error: %v\n", err)
		return
	}
	fmt.Printf("parseCSV(%q) = %q\n", csvLine, parsed)

	slugInput := "Hello World!"
	slug := tsk5.Slugify(slugInput)
	fmt.Printf("slugify(%q) = %q\n", slugInput, slug)

	logLine := "[INFO] 2026-08-04T12:00:00Z - Started"
	level, timestamp, message, err := tsk5.ParseLog(logLine)
	if err != nil {
		fmt.Printf("parseLog error: %v\n", err)
		return
	}
	fmt.Printf("parseLog(%q) = level=%q, timestamp=%q, message=%q\n", logLine, level, timestamp, message)
}
