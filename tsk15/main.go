package main

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store interface {
	GetAll() []Task
	GetByID(id int) (Task, error)
	Create(task Task) (Task, error)
	Update(id int, task Task) (Task, error)
	Delete(id int) error
}

type MemoryStore struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (s *MemoryStore) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	slices.SortFunc(result, func(a, b Task) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return result
}

func (s *MemoryStore) GetByID(id int) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task %d: %w", id, ErrNotFound)
	}
	return task, nil
}

func (s *MemoryStore) Create(task Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	task.CreatedAt = time.Now().UTC()
	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStore) Update(id int, task Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("update task %d: %w", id, ErrNotFound)
	}

	task.ID = id
	task.CreatedAt = current.CreatedAt
	s.tasks[id] = task
	return task, nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("delete task %d: %w", id, ErrNotFound)
	}
	delete(s.tasks, id)
	return nil
}

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("ошибка кодирования JSON: %v", err)
		}
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, value any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("decode JSON: multiple JSON values: %w", ErrBadRequest)
		}
		return fmt.Errorf("decode trailing JSON data: %w", err)
	}
	return nil
}

func parseID(value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id %q: %w", value, ErrBadRequest)
	}
	return id, nil
}

type TaskHandler struct {
	store Store
	mux   *http.ServeMux
}

func NewTaskHandler(store Store) *TaskHandler {
	handler := &TaskHandler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", handler.handleGetAll)
	mux.HandleFunc("POST /tasks", handler.handleCreate)
	mux.HandleFunc("GET /tasks/{id}", handler.handleGetByID)
	mux.HandleFunc("PUT /tasks/{id}", handler.handleUpdate)
	mux.HandleFunc("DELETE /tasks/{id}", handler.handleDelete)
	handler.mux = mux
	return handler
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *TaskHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.GetAll()

	values, filterProvided := r.URL.Query()["done"]
	if !filterProvided {
		writeJSON(w, http.StatusOK, tasks)
		return
	}
	if len(values) != 1 || (values[0] != "true" && values[0] != "false") {
		writeError(w, http.StatusBadRequest, "параметр done должен быть true или false")
		return
	}

	wantDone := values[0] == "true"
	filtered := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		if task.Done == wantDone {
			filtered = append(filtered, task)
		}
	}
	writeJSON(w, http.StatusOK, filtered)
}

func (h *TaskHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "задача не найдена")
			return
		}
		log.Printf("ошибка получения задачи %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := decodeJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON: "+err.Error())
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "поле title обязательно")
		return
	}

	created, err := h.store.Create(task)
	if err != nil {
		log.Printf("ошибка создания задачи: %v", err)
		writeError(w, http.StatusInternalServerError, "не удалось создать задачу")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var task Task
	if err := decodeJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON: "+err.Error())
		return
	}
	if strings.TrimSpace(task.Title) == "" {
		writeError(w, http.StatusBadRequest, "поле title обязательно")
		return
	}

	updated, err := h.store.Update(id, task)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "задача не найдена")
			return
		}
		log.Printf("ошибка обновления задачи %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "задача не найдена")
			return
		}
		log.Printf("ошибка удаления задачи %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s", r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
		log.Printf("← %s %s (%v)", r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func run() error {
	store := NewMemoryStore()
	if _, err := store.Create(Task{Title: "Изучить Go", Description: "Пройти курс"}); err != nil {
		return fmt.Errorf("create initial task: %w", err)
	}
	if _, err := store.Create(Task{Title: "Написать API", Description: "CRUD для задач"}); err != nil {
		return fmt.Errorf("create initial task: %w", err)
	}

	handler := loggingMiddleware(NewTaskHandler(store))
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Сервер запущен на http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case sig := <-quit:
		log.Printf("Получен сигнал %s, останавливаем сервер", sig)
	case err := <-serverErrors:
		return fmt.Errorf("serve HTTP: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	log.Println("Сервер остановлен")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
