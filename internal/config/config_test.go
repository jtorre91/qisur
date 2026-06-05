package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test_db")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")
	os.Setenv("PORT", "8080")
	os.Setenv("SEED", "false")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Config is nil")
	}

	if cfg.DatabaseURL != "postgres://test:test@localhost/test_db" {
		t.Errorf("DatabaseURL mismatch: got %s", cfg.DatabaseURL)
	}

	if cfg.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret mismatch: got %s", cfg.JWTSecret)
	}

	if cfg.JWTExpirationHours != 24 {
		t.Errorf("JWTExpirationHours mismatch: got %d, want 24", cfg.JWTExpirationHours)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port mismatch: got %s, want 8080", cfg.Port)
	}

	if cfg.Seed != false {
		t.Errorf("Seed mismatch: got %v, want false", cfg.Seed)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Test default values
	if cfg.Port != "8080" {
		t.Errorf("Default Port mismatch: got %s, want 8080", cfg.Port)
	}

	if cfg.JWTExpirationHours != 24 {
		t.Errorf("Default JWTExpirationHours mismatch: got %d, want 24", cfg.JWTExpirationHours)
	}

	if cfg.WSMaxMessageSize != 524288 {
		t.Errorf("Default WSMaxMessageSize mismatch: got %d, want 524288", cfg.WSMaxMessageSize)
	}

	if cfg.WSClientSendBuffer != 512 {
		t.Errorf("Default WSClientSendBuffer mismatch: got %d, want 512", cfg.WSClientSendBuffer)
	}

	if cfg.WSHubBroadcastBuffer != 1024 {
		t.Errorf("Default WSHubBroadcastBuffer mismatch: got %d, want 1024", cfg.WSHubBroadcastBuffer)
	}

	if cfg.WSWriteWait != 10*time.Second {
		t.Errorf("Default WSWriteWait mismatch: got %v, want 10s", cfg.WSWriteWait)
	}

	if cfg.WSPongWait != 60*time.Second {
		t.Errorf("Default WSPongWait mismatch: got %v, want 60s", cfg.WSPongWait)
	}
}

func TestLoadConfigWebSocketSettings(t *testing.T) {
	os.Setenv("WS_MAX_MESSAGE_SIZE", "1048576") // 1MB
	os.Setenv("WS_CLIENT_SEND_BUFFER", "256")
	os.Setenv("WS_HUB_BROADCAST_BUFFER", "512")
	os.Setenv("WS_WRITE_WAIT_SECONDS", "5")
	os.Setenv("WS_PONG_WAIT_SECONDS", "30")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.WSMaxMessageSize != 1048576 {
		t.Errorf("WSMaxMessageSize mismatch: got %d, want 1048576", cfg.WSMaxMessageSize)
	}

	if cfg.WSClientSendBuffer != 256 {
		t.Errorf("WSClientSendBuffer mismatch: got %d, want 256", cfg.WSClientSendBuffer)
	}

	if cfg.WSHubBroadcastBuffer != 512 {
		t.Errorf("WSHubBroadcastBuffer mismatch: got %d, want 512", cfg.WSHubBroadcastBuffer)
	}

	if cfg.WSWriteWait != 5*time.Second {
		t.Errorf("WSWriteWait mismatch: got %v, want 5s", cfg.WSWriteWait)
	}

	if cfg.WSPongWait != 30*time.Second {
		t.Errorf("WSPongWait mismatch: got %v, want 30s", cfg.WSPongWait)
	}
}
