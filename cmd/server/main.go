package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jtorre/qisurChallenge/internal/config"
	"github.com/jtorre/qisurChallenge/internal/db"
	"github.com/jtorre/qisurChallenge/internal/handlers"
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
	fmt.Println("✓ Server starting on port", cfg.Port)

	http.HandleFunc("/health", handlers.Health)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, http.DefaultServeMux))
}
