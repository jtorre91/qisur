package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jtorre/qisurChallenge/internal/config"
	"github.com/jtorre/qisurChallenge/internal/db"
	"github.com/jtorre/qisurChallenge/internal/router"
	"github.com/jtorre/qisurChallenge/seeds"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("✓ Connected to database")

	err = db.RunMigrations(pool)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("✓ Migrations completed")

	if cfg.Seed {
		if err = seeds.Run(pool); err != nil {
			log.Fatalf("failed to run seeders: %v", err)
		}
	} else {
		fmt.Println("✓ Seed skipped (set SEED=true to populate data)")
	}

	server := router.New(pool, cfg)

	fmt.Println("✓ Server starting on port", cfg.Port)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: server.Router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\n✓ Shutting down...")

		server.Hub.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		fmt.Println("✓ Server stopped")
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
