package v1

import (
	studentacademic "labs/api/v1/academics"
	bookingappointment "labs/api/v1/appointements"
	studentcontact "labs/api/v1/contact"
	studentdocument "labs/api/v1/document"
	studentenrollments "labs/api/v1/enrollement"
	"labs/api/v1/messages"
	"labs/api/v1/profilestudent"
	studentprofilesettings "labs/api/v1/settings"
	"labs/api/v1/universities"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesV1Init(router *gin.Engine, db *gorm.DB) {

	
	api  := router.Group("/v1")
	{
		
		profilestudent.StudentProfileRouterInit(api,db)
		studentcontact.StudentContactRouterInit(api,db)
		studentacademic.StudentAcademicRouterInit(api,db)
		studentdocument.StudentDocumentRouterInit(api,db)
		studentprofilesettings.StudentSettingsRouterInit(api,db)
		bookingappointment.AppointmentRouterInit(api,db)
		universities.UniversityRouterInit(api,db)
		studentenrollments.StudentEnrollmentRouterInit(api,db)
		messages.MessageRouterInit(api,db)
	}
}
