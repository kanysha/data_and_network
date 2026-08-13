package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func main() {
	targetURL := flag.String("url", "https://httpbin.org/get", "URL для GET-запроса")
	timeout := flag.Duration("timeout", 10*time.Second, "таймаут запроса")
	flag.Parse()

	client := &http.Client{Timeout: *timeout}
	resp, err := client.Get(*targetURL)
	if err != nil {
		log.Fatalf("выполнить GET %s: %v", *targetURL, err)
	}
	defer resp.Body.Close()

	fmt.Println("Статус:", resp.Status)
	fmt.Println("Код:", resp.StatusCode)
	fmt.Println("Протокол:", resp.Proto)
	fmt.Println()

	fmt.Println("Заголовки ответа:")
	keys := make([]string, 0, len(resp.Header))
	for key := range resp.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("  %s: %s\n", key, strings.Join(resp.Header.Values(key), ", "))
	}
	fmt.Println()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("прочитать тело ответа: %v", err)
	}
	fmt.Println("Тело ответа:")
	fmt.Println(string(body))
}
