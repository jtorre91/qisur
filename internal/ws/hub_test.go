package ws

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jtorre/qisurChallenge/internal/config"
)

func createTestConfig() *config.Config {
	return &config.Config{
		WSMaxMessageSize:     524288,
		WSClientSendBuffer:   512,
		WSHubBroadcastBuffer: 1024,
		WSWriteWait:          10 * time.Second,
		WSPongWait:           60 * time.Second,
	}
}

func TestNewHub(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)

	if hub == nil {
		t.Fatal("NewHub returned nil")
	}

	if hub.ClientCount() != 0 {
		t.Errorf("Initial ClientCount should be 0, got %d", hub.ClientCount())
	}
}

func TestHubBroadcast(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	type TestData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	data := TestData{
		ID:   "550e8400-e29b-41d4-a716-446655440000",
		Name: "Test Product",
	}

	err := hub.Broadcast("product_created", data)
	if err != nil {
		t.Fatalf("Broadcast failed: %v", err)
	}
}

func TestHubBroadcastInvalidData(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	// Create channel that can't be marshaled to JSON
	invalidData := make(chan struct{})

	err := hub.Broadcast("test_event", invalidData)
	if err == nil {
		t.Fatal("Expected error for invalid data")
	}
}

func TestClientCount(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	if hub.ClientCount() != 0 {
		t.Errorf("Initial ClientCount should be 0, got %d", hub.ClientCount())
	}

	client := &Client{
		id:   uuid.New(),
		hub:  hub,
		send: make(chan *Message, cfg.WSClientSendBuffer),
	}

	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond) // Wait for goroutine

	if hub.ClientCount() != 1 {
		t.Errorf("ClientCount after register should be 1, got %d", hub.ClientCount())
	}

	hub.UnregisterClient(client)
	time.Sleep(100 * time.Millisecond) // Wait for goroutine

	if hub.ClientCount() != 0 {
		t.Errorf("ClientCount after unregister should be 0, got %d", hub.ClientCount())
	}
}

func TestHubShutdown(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	client := &Client{
		id:   uuid.New(),
		hub:  hub,
		send: make(chan *Message, cfg.WSClientSendBuffer),
	}

	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	hub.Shutdown()
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("ClientCount after shutdown should be 0, got %d", hub.ClientCount())
	}
}

func TestRegisterUnregisterClient(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	client := &Client{
		id:   uuid.New(),
		hub:  hub,
		send: make(chan *Message, cfg.WSClientSendBuffer),
	}

	hub.RegisterClient(client)
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("Expected 1 client, got %d", hub.ClientCount())
	}

	hub.UnregisterClient(client)
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after unregister, got %d", hub.ClientCount())
	}
}

func TestMultipleClients(t *testing.T) {
	cfg := createTestConfig()
	hub := NewHub(cfg)
	go hub.Run()

	numClients := 5

	clients := make([]*Client, numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = &Client{
			id:   uuid.New(),
			hub:  hub,
			send: make(chan *Message, cfg.WSClientSendBuffer),
		}
		hub.RegisterClient(clients[i])
	}

	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != numClients {
		t.Errorf("Expected %d clients, got %d", numClients, hub.ClientCount())
	}

	// Unregister all clients
	for _, client := range clients {
		hub.UnregisterClient(client)
	}

	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after all unregistered, got %d", hub.ClientCount())
	}
}
