package studentacademic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentAcademicRouterInit initializes the routes related to student academic background.
func StudentAcademicRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewStudentAcademicRepository(db)

	// Private
	academic := router.Group("/student-academic")
	{
		// POST endpoint to create a new academic record
		academic.POST("", baseInstance.CreateOrUpdateStudentAcademic)

		// GET endpoint to retrieve the user's academic record
		academic.GET("", baseInstance.ReadStudentAcademic)

		// PUT endpoint to update the user's academic record
		academic.GET("/:userID", baseInstance.AdminGetStudentAcademicByID)

		// DELETE endpoint to delete the user's academic record
		academic.DELETE("", baseInstance.DeleteStudentAcademic)
	}
}
