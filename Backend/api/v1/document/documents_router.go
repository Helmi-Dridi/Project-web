package studentdocument

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentDocumentRouterInit initializes routes for student document handling.
func StudentDocumentRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize database wrapper
	baseInstance := Database{DB: db}

	// Auto-migrate if needed
	NewStudentDocumentRepository(db)

	// Routes group: /documents
	documents := router.Group("/documents")
	{
		// Core features
		documents.POST("/upload", baseInstance.UploadStudentDocument)
		documents.GET("", baseInstance.ReadStudentDocuments)
		documents.GET("/list", baseInstance.ReadStudentDocumentList)
		documents.GET("/:ID", baseInstance.ReadStudentDocument)
		documents.DELETE("/:ID", baseInstance.DeleteStudentDocument)

		// ✅ NEW endpoints for full UI support
		documents.GET("/status", baseInstance.GetStudentDocumentStatus)            // completion %
		documents.GET("/:ID/download", baseInstance.DownloadStudentDocument)       // download
		documents.POST("/:ID/finalize", baseInstance.FinalizeStudentDocument)     // mark as final
		documents.PUT("/:ID", baseInstance.UpdateStudentDocument) 
		documents.POST("/:ID/verify", baseInstance.VerifyStudentDocument)
		documents.GET("/user/:userID", baseInstance.AdminReadStudentDocuments)

                 // update document metadata or re-upload
	}

	
}
