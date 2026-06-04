package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
	Port               string
	Seed               bool
}

func Load() (*Config, error) {
	godotenv.Load()

	seed := false
	if os.Getenv("SEED") == "true" {
		seed = true
	}

	jwtHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if jwtHours == 0 {
		jwtHours = 24
	}

	return &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: jwtHours,
		Port:               getOrDefault("PORT", "8080"),
		Seed:               seed,
	}, nil
}

func getOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
