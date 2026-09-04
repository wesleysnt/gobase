package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/auth"
	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/httputil"
)

type Handler struct {
	svc *UserService
}

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	svc := &UserService{repo: &UserRepo{db: db}}
	h := &Handler{svc: svc}

	r.Group(func(r chi.Router) {
		r.Get("/users", h.List)
		r.Post("/users", h.Create)
		r.Get("/users/{id}", h.GetByID)
		r.Patch("/users/{id}", h.Update)
		r.Delete("/users/{id}", h.Delete)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth([]byte(cfg.JWTSecret)))
		r.Get("/me", h.Me)
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	if users == nil {
		users = []*UserResponse{}
	}

	httputil.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil {
		httputil.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := h.svc.GetByID(r.Context(), claims.Subject)
	if err != nil {
		httputil.WriteError(w, httputil.CodeFrom(err), err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}
