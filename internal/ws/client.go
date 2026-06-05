package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/jtorre/qisurChallenge/internal/config"
)

type Client struct {
	id     uuid.UUID
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Message
	config *config.Config
}

func NewClient(id uuid.UUID, hub *Hub, conn *websocket.Conn, cfg *config.Config) *Client {
	return &Client{
		id:     id,
		hub:    hub,
		conn:   conn,
		send:   make(chan *Message, cfg.WSClientSendBuffer),
		config: cfg,
	}
}

func (c *Client) Run() {
	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(int64(c.config.WSMaxMessageSize))
	c.conn.SetReadDeadline(time.Now().Add(c.config.WSPongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.config.WSPongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
			}
			return
		}
	}
}

func (c *Client) writePump() {
	pingInterval := (c.config.WSPongWait * 9) / 10
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.config.WSWriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.config.WSWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
