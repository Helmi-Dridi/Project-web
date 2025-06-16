package universities

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UniversityRouterInit initializes routes for universities and their programs.
func UniversityRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	// Create database instance
	baseInstance := Database{DB: db}

	// Perform AutoMigrate
	NewUniversityRepository(db)

	// Group routes under /universities
	universities := router.Group("/universities")
	{
		// GET /universities/:companyID — get all universities for a company
		universities.GET("/:companyID", baseInstance.GetAllUniversities)

		// GET /universities/:companyID/:id/programs — get programs of a specific university
		universities.GET("/:companyID/:id/programs", baseInstance.GetUniversityPrograms)
	}
}
