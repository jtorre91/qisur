package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jtorre/qisurChallenge/internal/config"
	"github.com/jtorre/qisurChallenge/internal/handlers"
	"github.com/jtorre/qisurChallenge/internal/middleware"
	"github.com/jtorre/qisurChallenge/internal/repository"
	"github.com/jtorre/qisurChallenge/internal/ws"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/jtorre/qisurChallenge/docs"
)

func New(pool *pgxpool.Pool, cfg *config.Config) chi.Router {
	// Initialize WebSocket
	hub := ws.NewHub()
	go hub.Run()

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	searchRepo := repository.NewSearchRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo, hub)
	productHandler := handlers.NewProductHandler(productRepo, hub)
	searchHandler := handlers.NewSearchHandler(searchRepo)
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	wsHandler := handlers.NewWSHandler(hub, cfg)

	// Setup router
	router := chi.NewRouter()

	// Health check
	router.Get("/health", handlers.Health)

	// Swagger documentation
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Auth routes
	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Categories routes
	router.Route("/api/categories", func(r chi.Router) {
		r.Get("/", categoryHandler.List)
		r.Get("/{id}", categoryHandler.GetByID)

		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Post("/", categoryHandler.Create)
		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Put("/{id}", categoryHandler.Update)
		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Delete("/{id}", categoryHandler.Delete)
	})

	// Products routes
	router.Route("/api/products", func(r chi.Router) {
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.GetByID)
		r.Get("/{id}/history", productHandler.GetHistory)

		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Post("/", productHandler.Create)
		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Put("/{id}", productHandler.Update)
		r.With(middleware.AuthMiddleware(cfg), middleware.RoleGuard("admin")).Delete("/{id}", productHandler.Delete)
	})

	// Search routes
	router.Get("/api/search", searchHandler.Search)

	// WebSocket route
	router.Get("/ws", wsHandler.Handle)

	return router
}
