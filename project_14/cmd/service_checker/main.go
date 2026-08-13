package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type CheckResult struct {
	URL    string
	Status int
	Err    error
	Took   time.Duration
}

func checkURL(ctx context.Context, client *http.Client, targetURL string) CheckResult {
	start := time.Now()
	result := CheckResult{URL: targetURL}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		result.Err = fmt.Errorf("создать запрос: %w", err)
		result.Took = time.Since(start)
		return result
	}

	resp, err := client.Do(req)
	result.Took = time.Since(start)
	if err != nil {
		result.Err = fmt.Errorf("выполнить запрос: %w", err)
		return result
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64*1024))
	result.Status = resp.StatusCode
	return result
}

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		urls = []string{
			"https://google.com",
			"https://github.com",
			"https://golang.org",
			"https://httpbin.org/status/500",
			"https://nonexistent.invalid",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := &http.Client{Timeout: 5 * time.Second}
	results := make(chan CheckResult, len(urls))

	var wg sync.WaitGroup
	for _, targetURL := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- checkURL(ctx, client, targetURL)
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			fmt.Printf("[FAIL] %s — %v (%v)\n", result.URL, result.Err, result.Took.Round(time.Millisecond))
			continue
		}
		fmt.Printf("[%d] %s — %v\n", result.Status, result.URL, result.Took.Round(time.Millisecond))
	}
}
