package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jtorre/qisurChallenge/internal/auth"
	"github.com/jtorre/qisurChallenge/internal/config"
	"github.com/jtorre/qisurChallenge/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *ws.Hub
	cfg *config.Config
}

func NewWSHandler(hub *ws.Hub, cfg *config.Config) *WSHandler {
	return &WSHandler{
		hub: hub,
		cfg: cfg,
	}
}

func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token, h.cfg.JWTSecret)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	client := ws.NewClient(claims.UserID, h.hub, conn, h.cfg)
	h.hub.RegisterClient(client)

	client.Run()
}
