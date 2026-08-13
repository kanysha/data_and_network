package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func corsFileHandler(directory string) http.Handler {
	files := http.FileServer(http.Dir(directory))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}
		files.ServeHTTP(w, r)
	})
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP-адрес для прослушивания")
	directory := flag.String("dir", "project_12/public", "каталог со статическими файлами")
	flag.Parse()

	info, err := os.Stat(*directory)
	if err != nil {
		log.Fatalf("проверить каталог %s: %v", *directory, err)
	}
	if !info.IsDir() {
		log.Fatalf("%s не является каталогом", *directory)
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           corsFileHandler(*directory),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Файловый сервер запущен на http://localhost%s", *addr)
	log.Printf("Директория: %s", *directory)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
