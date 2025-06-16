package studentcontact

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentContactRouterInit initializes the routes related to student contact details.
func StudentContactRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewStudentContactRepository(db)

	// Private
	contact := router.Group("/student-contact")
	{
		// POST endpoint to create a new student contact record
		contact.POST("", baseInstance.CreateOrUpdateStudentContact)

		// GET endpoint to retrieve the student's contact information
		contact.GET("", baseInstance.ReadStudentContact)

		// PUT endpoint to update the student's contact information
		contact.GET("/:userID", baseInstance.AdminGetStudentContactByID)

		// DELETE endpoint to delete the student's contact information
		contact.DELETE("", baseInstance.DeleteStudentContact)
	}
}
