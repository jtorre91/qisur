package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL           string
	JWTSecret             string
	JWTExpirationHours    int
	Port                  string
	Seed                  bool
	WSMaxMessageSize      int
	WSClientSendBuffer    int
	WSHubBroadcastBuffer  int
	WSWriteWait           time.Duration
	WSPongWait            time.Duration
	WSPingInterval        time.Duration
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

	// WebSocket configuration
	wsMaxMsgSize, _ := strconv.Atoi(getOrDefault("WS_MAX_MESSAGE_SIZE", "524288")) // 512KB
	wsClientBuffer, _ := strconv.Atoi(getOrDefault("WS_CLIENT_SEND_BUFFER", "512"))
	wsHubBuffer, _ := strconv.Atoi(getOrDefault("WS_HUB_BROADCAST_BUFFER", "1024"))

	writeWaitSecs, _ := strconv.Atoi(getOrDefault("WS_WRITE_WAIT_SECONDS", "10"))
	pongWaitSecs, _ := strconv.Atoi(getOrDefault("WS_PONG_WAIT_SECONDS", "60"))

	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		JWTExpirationHours:   jwtHours,
		Port:                 getOrDefault("PORT", "8080"),
		Seed:                 seed,
		WSMaxMessageSize:     wsMaxMsgSize,
		WSClientSendBuffer:   wsClientBuffer,
		WSHubBroadcastBuffer: wsHubBuffer,
		WSWriteWait:          time.Duration(writeWaitSecs) * time.Second,
		WSPongWait:           time.Duration(pongWaitSecs) * time.Second,
		WSPingInterval:       0, // Calculated as (pongWait * 9) / 10
	}, nil
}

func getOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
