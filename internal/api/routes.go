package api

import (
	"net/http"
	"strade/internal/api/handlers"
	stradeMiddleware "strade/internal/api/middleware"
	"strade/internal/env"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func NewRouter(h *handlers.HandlerManager, bundle *i18n.Bundle) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(stradeMiddleware.Localizer(bundle))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   env.GetSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.HandleHealthCheck)
	})

	return r
}
