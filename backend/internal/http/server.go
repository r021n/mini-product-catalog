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

	userStore := store.NewUserStore(db)
	categoryStore := store.NewCategoryStore(db)
	productStore := store.NewProductStore(db)

	healthHandler := handler.NewHealthHandler()
	categoriesHandler := handler.NewCategoriesHandler(categoryStore, validate)
	productsHandler := handler.NewProductsHandler(productStore, categoryStore, validate)
	authHandler := handler.NewAuthHandler(userStore, validate, cfg.JWTSecret)

	r.Get("/health", healthHandler.Health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/me", authHandler.Me)
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", categoriesHandler.List)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", categoriesHandler.Create)
			r.Put("/{id}", categoriesHandler.Update)
			r.Delete("/{id}", categoriesHandler.Delete)
		})
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productsHandler.List)
		r.Get("/{id}", productsHandler.Get)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", productsHandler.Create)
			r.Put("/{id}", productsHandler.Update)
			r.Delete("/{id}", productsHandler.Delete)
		})
	})

	return r
}
