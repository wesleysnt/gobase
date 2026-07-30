package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupTestRouter() chi.Router {
	r := chi.NewRouter()
	repo := newStubRepo()
	svc := &UserService{repo: repo}
	h := &Handler{svc: svc}

	r.Post("/users", h.Create)
	r.Get("/users", h.List)
	r.Get("/users/{id}", h.GetByID)
	r.Patch("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)

	return r
}

func TestHandlerCreate(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"Alice","email":"alice@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp UserResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", resp.Name)
	}
	if resp.ID == "" {
		t.Error("ID is empty")
	}
}

func TestHandlerCreateValidation(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"","email":""}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestHandlerGetByID(t *testing.T) {
	router := setupTestRouter()

	// Create a user first
	createBody := `{"name":"Bob","email":"bob@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var created UserResponse
	json.NewDecoder(rec.Body).Decode(&created)

	// Now fetch it
	req = httptest.NewRequest("GET", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var fetched UserResponse
	json.NewDecoder(rec.Body).Decode(&fetched)
	if fetched.Name != "Bob" {
		t.Errorf("Name = %q, want Bob", fetched.Name)
	}
}

func TestHandlerGetByIDNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/users/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandlerList(t *testing.T) {
	router := setupTestRouter()

	// Create two users
	for _, body := range []string{
		`{"name":"Alice","email":"alice@example.com"}`,
		`{"name":"Bob","email":"bob@example.com"}`,
	} {
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest("GET", "/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var users []UserResponse
	json.NewDecoder(rec.Body).Decode(&users)
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}
}

func TestHandlerDelete(t *testing.T) {
	router := setupTestRouter()

	// Create a user
	createBody := `{"name":"Charlie","email":"charlie@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var created UserResponse
	json.NewDecoder(rec.Body).Decode(&created)

	// Delete it
	req = httptest.NewRequest("DELETE", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's gone
	req = httptest.NewRequest("GET", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d (not found after delete)", rec.Code, http.StatusNotFound)
	}
}

func TestHandlerUpdate(t *testing.T) {
	router := setupTestRouter()

	// Create a user
	createBody := `{"name":"Diana","email":"diana@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var created UserResponse
	json.NewDecoder(rec.Body).Decode(&created)

	// Update the user
	updateBody := `{"name":"Diana Updated"}`
	req = httptest.NewRequest("PATCH", "/users/"+created.ID, strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var updated UserResponse
	json.NewDecoder(rec.Body).Decode(&updated)
	if updated.Name != "Diana Updated" {
		t.Errorf("Name = %q, want Diana Updated", updated.Name)
	}
}

func TestHandlerUpdateNotFound(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"Nobody"}`
	req := httptest.NewRequest("PATCH", "/users/nonexistent", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
