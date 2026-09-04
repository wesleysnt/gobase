package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/user"
)

func NewRouter(cfg *config.Config, db *sqlx.DB, log *slog.Logger) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(RequestLogger(log))
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 — domain routes registered here
	r.Route("/api/v1", func(r chi.Router) {
		user.RegisterRoutes(r, db, cfg)
	})

	return r
}
