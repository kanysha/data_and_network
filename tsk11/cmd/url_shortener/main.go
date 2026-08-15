package main

import (
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const codeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLStore struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func (s *URLStore) Save(code, original string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.urls[code]; exists {
		return false
	}
	s.urls[code] = original
	return true
}

func (s *URLStore) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	original, ok := s.urls[code]
	return original, ok
}

func (s *URLStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func generateCode(length int) (string, error) {
	code := make([]byte, length)
	limit := big.NewInt(int64(len(codeAlphabet)))
	for i := range code {
		index, err := cryptorand.Int(cryptorand.Reader, limit)
		if err != nil {
			return "", fmt.Errorf("generate random code: %w", err)
		}
		code[i] = codeAlphabet[index.Int64()]
	}
	return string(code), nil
}

type Shortener struct {
	store   *URLStore
	baseURL string
}

func NewShortener(store *URLStore, baseURL string) *Shortener {
	return &Shortener{store: store, baseURL: strings.TrimSuffix(baseURL, "/")}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode JSON response: %v", err)
	}
}

func (s *Shortener) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "используйте POST"})
		return
	}
	defer r.Body.Close()

	var request struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "неверный JSON"})
		return
	}
	if err := validateURL(request.URL); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var code string
	for range 10 {
		generated, err := generateCode(6)
		if err != nil {
			log.Printf("generate code: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "не удалось создать короткую ссылку"})
			return
		}
		if s.store.Save(generated, request.URL) {
			code = generated
			break
		}
	}
	if code == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "не удалось подобрать уникальный код"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"code":  code,
		"short": s.baseURL + "/s/" + code,
	})
}

func validateURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("url должен быть абсолютным HTTP- или HTTPS-адресом")
	}
	return nil
}

func (s *Shortener) redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "используйте GET"})
		return
	}
	code := strings.TrimPrefix(r.URL.Path, "/s/")
	if code == "" || strings.Contains(code, "/") {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "код не указан"})
		return
	}

	original, ok := s.store.Get(code)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "ссылка не найдена"})
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

func (s *Shortener) statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "используйте GET"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"total": s.store.Len()})
}

func main() {
	shortener := NewShortener(NewURLStore(), "http://localhost:8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortener.shortenHandler)
	mux.HandleFunc("/s/", shortener.redirectHandler)
	mux.HandleFunc("/stats", shortener.statsHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("URL Shortener запущен на http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
