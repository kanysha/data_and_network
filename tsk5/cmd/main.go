package main

import (
	"fmt"

	project5 "github.com/kanysha/data_and_network/project_5"
)

func main() {
	csvLine := "1,\"two, too\",3"
	parsed := project5.ParseCSV(csvLine)
	fmt.Printf("parseCSV(%q) = %q\n", csvLine, parsed)

	slugInput := "Hello World!"
	slug := project5.Slugify(slugInput)
	fmt.Printf("slugify(%q) = %q\n", slugInput, slug)

	logLine := "[INFO] 2026-08-04T12:00:00Z - Started"
	level, timestamp, message, err := project5.ParseLog(logLine)
	if err != nil {
		fmt.Printf("parseLog error: %v\n", err)
		return
	}
	fmt.Printf("parseLog(%q) = level=%q, timestamp=%q, message=%q\n", logLine, level, timestamp, message)
}
