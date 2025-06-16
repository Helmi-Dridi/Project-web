package studentenrollments

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StudentEnrollmentRouterInit(router *gin.RouterGroup, db *gorm.DB) {
	base := Database{DB: db}
	// Automigrate / Update table
	NewStudentEnrollmentRepository(db)
	enrollments := router.Group("/enrollments/:companyID")
	{
		enrollments.POST("/users/:userID/programs/:programID", base.AssignStudentToProgram)
		enrollments.DELETE("/users/:userID", base.UnassignStudentFromProgram)
		enrollments.GET("/view", base.GetStudentEnrollment)
		enrollments.GET("/programs/:programID", base.GetEnrollmentsByProgramID)
		enrollments.GET("/view/all", base.GetAllStudentEnrollments)
	}

}
