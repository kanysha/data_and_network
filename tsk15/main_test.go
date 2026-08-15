package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGetAllAndFilter(t *testing.T) {
	store := NewMemoryStore()
	if _, err := store.Create(Task{Title: "Не завершена", Done: false}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Create(Task{Title: "Завершена", Done: true}); err != nil {
		t.Fatal(err)
	}
	handler := NewTaskHandler(store)

	tests := []struct {
		name       string
		target     string
		wantStatus int
		wantCount  int
		wantDone   *bool
	}{
		{name: "all", target: "/tasks", wantStatus: http.StatusOK, wantCount: 2},
		{name: "done", target: "/tasks?done=true", wantStatus: http.StatusOK, wantCount: 1, wantDone: boolPtr(true)},
		{name: "not done", target: "/tasks?done=false", wantStatus: http.StatusOK, wantCount: 1, wantDone: boolPtr(false)},
		{name: "invalid filter", target: "/tasks?done=unknown", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.target, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantStatus != http.StatusOK {
				return
			}

			var tasks []Task
			if err := json.NewDecoder(rec.Body).Decode(&tasks); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if len(tasks) != tt.wantCount {
				t.Fatalf("task count = %d, want %d", len(tasks), tt.wantCount)
			}
			if tt.wantDone != nil && tasks[0].Done != *tt.wantDone {
				t.Fatalf("task done = %t, want %t", tasks[0].Done, *tt.wantDone)
			}
		})
	}
}

func TestCRUDStatuses(t *testing.T) {
	handler := NewTaskHandler(NewMemoryStore())

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Тест"}`))
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want %d", createRec.Code, http.StatusCreated)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /tasks/1 status = %d, want %d", getRec.Code, http.StatusOK)
	}

	updateReq := httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(`{"title":"Обновлено","done":true}`))
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("PUT /tasks/1 status = %d, want %d; body: %s", updateRec.Code, http.StatusOK, updateRec.Body.String())
	}

	var updated Task
	if err := json.NewDecoder(updateRec.Body).Decode(&updated); err != nil {
		t.Fatalf("decode PUT response: %v", err)
	}
	if updated.Title != "Обновлено" || !updated.Done {
		t.Fatalf("PUT response = %+v, want updated title and done=true", updated)
	}

	missingReq := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("GET missing task status = %d, want %d", missingRec.Code, http.StatusNotFound)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/1 status = %d, want %d", deleteRec.Code, http.StatusNoContent)
	}

	deletedReq := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	deletedRec := httptest.NewRecorder()
	handler.ServeHTTP(deletedRec, deletedReq)
	if deletedRec.Code != http.StatusNotFound {
		t.Fatalf("GET deleted task status = %d, want %d", deletedRec.Code, http.StatusNotFound)
	}
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty title", body: `{"title":"  "}`},
		{name: "malformed JSON", body: `{"title":`},
		{name: "unknown field", body: `{"title":"Тест","unknown":true}`},
		{name: "multiple values", body: `{"title":"Тест"} {"title":"Лишнее"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTaskHandler(NewMemoryStore())
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("POST /tasks status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
		})
	}
}

func TestRouting(t *testing.T) {
	handler := NewTaskHandler(NewMemoryStore())
	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{name: "extra path segment", method: http.MethodGet, target: "/tasks/1/extra", wantStatus: http.StatusNotFound},
		{name: "unsupported collection method", method: http.MethodDelete, target: "/tasks", wantStatus: http.StatusMethodNotAllowed},
		{name: "invalid id", method: http.MethodGet, target: "/tasks/not-a-number", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("%s %s status = %d, want %d", tt.method, tt.target, rec.Code, tt.wantStatus)
			}
		})
	}
}

func boolPtr(value bool) *bool {
	return &value
}
