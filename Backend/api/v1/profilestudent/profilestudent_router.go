package profilestudent

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentProfileRouterInit initializes the routes related to student profiles.
func StudentProfileRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewStudentProfileRepository(db)

	// Private
	profile := router.Group("/student-profile")
	{
		// POST endpoint to create a new student profile
		profile.POST("", baseInstance.CreateOrUpdateStudentProfile)

		// GET endpoint to retrieve the student profile of the current user
		profile.GET("", baseInstance.ReadStudentProfile)

		
		profile.GET("/:userID", baseInstance.AdminGetStudentProfileByID)

		// DELETE endpoint to delete the student profile of the current user
		profile.DELETE("", baseInstance.DeleteStudentProfile)
	}
}
