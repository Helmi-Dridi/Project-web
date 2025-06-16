package messages

import (
	"labs/domains"
	"labs/utils"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

type PresenceUpdate struct {
	Type   string    `json:"type"`   // "presence"
	UserID uuid.UUID `json:"user_id"`
	Status string    `json:"status"` // "online" | "offline"
}

type IncomingMessage struct {
	Type       string  `json:"type"`       // "message" | "typing"
	To         string  `json:"to"`         // UUID string
	Content    string  `json:"content"`    // Only for message
	Attachment *string `json:"attachment"` // Optional
}

// ✅ Outgoing message struct with all required fields
type WebSocketMessage struct {
	Type       string     `json:"type"`
	ID         uuid.UUID  `json:"id"`
	From       string     `json:"from"`
	SenderID   uuid.UUID  `json:"sender_id"`
	ReceiverID uuid.UUID  `json:"receiver_id"`
	CompanyID  uuid.UUID  `json:"company_id"`
	Content    string     `json:"content"`
	Attachment *string    `json:"attachment,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

var (
	clients   = make(map[uuid.UUID]*Client)
	clientsMu sync.RWMutex
)

func (db Database) WebSocketHandler(ctx *gin.Context) {
	if token := ctx.Query("token"); token != "" {
		ctx.Request.Header.Set("Authorization", "Bearer "+token)
	}

	session := utils.ExtractJWTValues(ctx)
	if session.UserID == uuid.Nil || session.CompanyID == uuid.Nil {
		log.Println("WebSocket rejected: invalid or missing token claims")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	userID := session.UserID
	companyID := session.CompanyID

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	clientsMu.Lock()
	clients[userID] = &Client{UserID: userID, Conn: conn}
	clientsMu.Unlock()

	broadcastPresence(userID, "online")
	go pingClient(conn)

	defer func() {
		conn.Close()
		clientsMu.Lock()
		delete(clients, userID)
		clientsMu.Unlock()
		broadcastPresence(userID, "offline")
	}()

	for {
		var msg IncomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		switch msg.Type {
		case "typing":
			targetID, err := uuid.Parse(msg.To)
			if err != nil {
				continue
			}
			clientsMu.RLock()
			receiver, ok := clients[targetID]
			clientsMu.RUnlock()
			if ok {
				_ = receiver.Conn.WriteJSON(gin.H{
					"type": "typing",
					"from": userID.String(),
				})
			}

		case "message":
			receiverID, err := uuid.Parse(msg.To)
			if err != nil || msg.Content == "" {
				continue
			}

			message := &domains.Message{
				ID:         uuid.New(),
				SenderID:   userID,
				ReceiverID: receiverID,
				Content:    msg.Content,
				Attachment: msg.Attachment,
				CompanyID:  companyID,
				CreatedAt:  time.Now(),
				Read:       false,
			}

			if err := CreateMessage(db.DB, message); err != nil {
				log.Println("Failed to save message:", err)
				continue
			}

			clientsMu.RLock()
			receiver, ok := clients[receiverID]
			clientsMu.RUnlock()
			if ok {
				outgoing := WebSocketMessage{
					Type:       "message",
					ID:         message.ID,
					From:       userID.String(),
					SenderID:   message.SenderID,
					ReceiverID: message.ReceiverID,
					CompanyID:  message.CompanyID,
					Content:    message.Content,
					Attachment: message.Attachment,
					CreatedAt:  message.CreatedAt,
				}
				_ = receiver.Conn.WriteJSON(outgoing)
			}
		}
	}
}

func pingClient(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			break
		}
	}
}

func broadcastPresence(userID uuid.UUID, status string) {
	update := PresenceUpdate{
		Type:   "presence",
		UserID: userID,
		Status: status,
	}

	clientsMu.RLock()
	defer clientsMu.RUnlock()
	for _, client := range clients {
		if client.UserID != userID {
			_ = client.Conn.WriteJSON(update)
		}
	}
}
