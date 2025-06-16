package messages

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MessageRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	base := Database{DB: db}
	NewMessageRepository(db)

	

	// Message API routes (scoped by company)
	messages := router.Group("/messages/:companyID")
	{
		messages.POST("/send", base.SendMessage)
		messages.GET("/user/:userID", base.GetConversation)
		messages.POST("/:id/read", base.MarkAsRead)
		messages.GET("/partners", base.GetConversationPartners)
		messages.GET("/unread-count", base.GetUnreadMessageCount)
		messages.GET("/user/:userID/paginated", base.GetPaginatedConversation)
		messages.DELETE("/:id", base.DeleteMessage)
		messages.GET("/search", base.SearchMessages)
		messages.GET("/inbox", base.GetAdminInbox)
		// WebSocket route (global)
		messages.GET("/ws", base.WebSocketHandler)
	}
}
