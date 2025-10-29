package api

import (
	"net/http"
	"strade/internal/api/handle"
	smw "strade/internal/api/middleware"
	"strade/internal/env"
	"strade/internal/i18n"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(h *handle.Handler) http.Handler {
	r := chi.NewRouter()
	bundle := i18n.NewBundle()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   env.GetSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(smw.Localizer(bundle))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.HandleHealthCheck)
	})

	return r
}
