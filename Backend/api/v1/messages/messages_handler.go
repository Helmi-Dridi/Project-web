package messages

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendMessageRequest struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Content    string    `json:"content"`
	Attachment *string   `json:"attachment,omitempty"`
}

// SendMessage handles sending a message.
// SendMessage handles sending a message.
func (db Database) SendMessage(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))

	var req SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	message := &domains.Message{
		ID:         uuid.New(),
		SenderID:   session.UserID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		Attachment: req.Attachment,
		CompanyID:  companyID,
		CreatedAt:  time.Now(),
		Read:       false,
	}

	if err := CreateMessage(db.DB, message); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to send message", utils.Null())
		return
	}

	// 🔄 WebSocket push to receiver if connected
	clientsMu.RLock()
	receiverClient, ok := clients[req.ReceiverID]
	clientsMu.RUnlock()

	if ok {
	outgoing := WebSocketMessage{
		Type:       "message",
		ID:         message.ID,
		From:       session.UserID.String(),
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		CompanyID:  message.CompanyID,
		Content:    message.Content,
		Attachment: message.Attachment,
		CreatedAt:  message.CreatedAt,
	}
	_ = receiverClient.Conn.WriteJSON(outgoing)
}

	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, message)
}


// GetConversation fetches all messages between logged-in user and another user.
func (db Database) GetConversation(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	userID, _ := uuid.Parse(ctx.Param("userID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	messages, err := GetMessagesBetween(db.DB, session.UserID, userID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch messages", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, messages)
}

// MarkAsRead sets a message's read status to true.
func (db Database) MarkAsRead(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	messageID, _ := uuid.Parse(ctx.Param("id"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	if err := MarkMessageRead(db.DB, messageID); err != nil {
	utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to update message", utils.Null())
	return
}

// Optional: Fetch the message to get sender ID
var msg domains.Message
if err := db.DB.First(&msg, "id = ?", messageID).Error; err == nil {
	clientsMu.RLock()
	senderClient, ok := clients[msg.SenderID]
	clientsMu.RUnlock()

	if ok {
		_ = senderClient.Conn.WriteJSON(gin.H{
			"type":       "read_receipt",
			"message_id": msg.ID,
			"reader_id":  session.UserID,
		})
	}
}


	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, "Message marked as read")
}
func (db Database) GetConversationPartners(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	partners, err := GetConversationPartners(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch conversation partners", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, partners)
}
func (db Database) GetUnreadMessageCount(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	count, err := CountUnreadMessages(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to count unread messages", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{"unreadCount": count})
}
func (db Database) GetPaginatedConversation(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	userID, _ := uuid.Parse(ctx.Param("userID"))

	// Parse limit & offset safely
	limitStr := ctx.DefaultQuery("limit", "20")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	messages, err := GetPaginatedMessages(db.DB, session.UserID, userID, limit, offset)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch paginated messages", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, messages)
}
func (db Database) DeleteMessage(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	messageID, _ := uuid.Parse(ctx.Param("id"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	if err := DeleteMessage(db.DB, messageID, session.UserID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete message", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, "Message deleted")
}
func (db Database) SearchMessages(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	query := ctx.Query("q")

	if query == "" {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Missing search query", utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	results, err := SearchMessages(db.DB, companyID, session.UserID, query)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to search messages", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, results)
}
func (db Database) GetAdminInbox(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	inbox, err := GetAdminInbox(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch admin inbox", utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, inbox)
}
