package http

import (
	"log/slog"
	"mini-product-catalog/internal/config"
	"mini-product-catalog/internal/http/handler"
	"mini-product-catalog/internal/middleware"
	"mini-product-catalog/internal/store"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewServer(cfg config.Config, logger *slog.Logger, db *pgxpool.Pool) nethttp.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)

	r.Use(middleware.CORS(cfg.AllowedOrigins))

	r.Use(middleware.RequestLogger(logger))

	validate := validator.New()

	categoryStore := store.NewCategoryStore(db)

	healthHandler := handler.NewHealthHandler()
	categoriesHandler := handler.NewCategoriesHandler(categoryStore, validate)

	r.Get("/health", healthHandler.Health)

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", categoriesHandler.List)
		r.Post("/", categoriesHandler.Create)
	})

	return r
}
