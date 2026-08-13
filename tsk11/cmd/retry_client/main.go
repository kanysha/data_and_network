package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

func doWithRetry(ctx context.Context, client *http.Client, url string, maxRetries int) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			fmt.Printf("Попытка %d через %v...\n", attempt+1, delay)
			if err := waitForRetry(ctx, delay); err != nil {
				return nil, err
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("создать запрос: %w", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			lastErr = fmt.Errorf("выполнить запрос: %w", err)
			fmt.Printf("Попытка %d завершилась ошибкой: %v\n", attempt+1, err)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("прочитать тело ответа: %w", readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("закрыть тело ответа: %w", closeErr)
		}
		if resp.StatusCode >= http.StatusInternalServerError {
			lastErr = fmt.Errorf("сервер вернул %s: %s", resp.Status, body)
			fmt.Printf("Попытка %d: сервер вернул %s\n", attempt+1, resp.Status)
			continue
		}
		return body, nil
	}

	if lastErr == nil {
		lastErr = errors.New("неизвестная ошибка")
	}
	return nil, fmt.Errorf("все %d попыток исчерпаны: %w", maxRetries+1, lastErr)
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func main() {
	targetURL := flag.String("url", "https://httpbin.org/get", "URL для GET-запроса")
	maxRetries := flag.Int("retries", 3, "количество повторных попыток")
	timeout := flag.Duration("timeout", 30*time.Second, "общий таймаут")
	flag.Parse()
	if *maxRetries < 0 {
		fmt.Println("Ошибка: retries не может быть отрицательным")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	client := &http.Client{Timeout: 5 * time.Second}
	body, err := doWithRetry(ctx, client, *targetURL, *maxRetries)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Успех! Получено %d байт\n%s\n", len(body), body)
}
