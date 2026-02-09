package httpapi

import (
	"log/slog"
	"net/http"

	"mini-product-catalog/internal/config"
	"mini-product-catalog/internal/httpapi/handlers"
	mw "mini-product-catalog/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(cfg config.Config, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(mw.CORS(cfg.AllowedOrigins))

	r.Use(mw.RequestLogger(logger))

	r.Get("/health", handlers.Health())

	return r
}
