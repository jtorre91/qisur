package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jtorre/qisurChallenge/internal/handlers"
	"github.com/jtorre/qisurChallenge/internal/repository"
)

func New(pool *pgxpool.Pool) chi.Router {
	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	searchRepo := repository.NewSearchRepository(pool)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	productHandler := handlers.NewProductHandler(productRepo)
	searchHandler := handlers.NewSearchHandler(searchRepo)

	// Setup router
	router := chi.NewRouter()

	// Health check
	router.Get("/health", handlers.Health)

	// Categories routes
	router.Route("/api/categories", func(r chi.Router) {
		r.Get("/", categoryHandler.List)
		r.Post("/", categoryHandler.Create)
		r.Get("/{id}", categoryHandler.GetByID)
		r.Put("/{id}", categoryHandler.Update)
		r.Delete("/{id}", categoryHandler.Delete)
	})

	// Products routes
	router.Route("/api/products", func(r chi.Router) {
		r.Get("/", productHandler.List)
		r.Post("/", productHandler.Create)
		r.Get("/{id}", productHandler.GetByID)
		r.Put("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)
		r.Get("/{id}/history", productHandler.GetHistory)
	})

	// Search routes
	router.Get("/api/search", searchHandler.Search)

	return router
}
