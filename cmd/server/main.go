package main

import (
	"fmt"
	"log"
	"net/http"

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

	r := router.New(pool)

	fmt.Println("✓ Server starting on port", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
