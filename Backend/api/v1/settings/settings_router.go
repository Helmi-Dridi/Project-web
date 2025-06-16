package studentprofilesettings

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentSettingsRouterInit initializes student settings routes
func StudentSettingsRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate table
	NewStudentSettingsRepository(db)

	// Group routes under /student/settings
	settings := router.Group("/student/settings")
	{
		settings.POST("", baseInstance.CreateOrUpdateStudentSettings) // Create or update settings
		settings.GET("", baseInstance.GetStudentSettings)             // Read settings
		settings.DELETE("", baseInstance.DeleteStudentSettingsHandler) // Delete settings
		settings.GET("/:userID", baseInstance.AdminGetStudentSettingsByID)
	}
}
