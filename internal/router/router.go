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

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)

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

	return router
}
