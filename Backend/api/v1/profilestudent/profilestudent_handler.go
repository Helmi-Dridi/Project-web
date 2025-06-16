package profilestudent

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateStudentProfile 	Handles the creation of a new student profile.
// @Summary         	Create student profile
// @Description     	Create a new student profile.
// @Tags				StudentProfile
// @Accept				json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param				request	body			studentprofile.StudentProfileIn	true	"Student profile input"
// @Success			201			{object}		utils.ApiResponses
// @Failure			400			{object}		utils.ApiResponses	"Invalid request"
// @Failure			401			{object}		utils.ApiResponses	"Unauthorized"
// @Failure			500			{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/student-profile	[post]
func (db Database) CreateStudentProfile(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(domains.StudentProfile)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid input. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	input.ID = uuid.New()
	input.UserID = session.UserID

	if err := db.DB.Create(input).Error; err != nil {
		logrus.Error("Failed to save student profile. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadStudentProfile 	Handles the retrieval of a student profile by UserID.
// @Summary         	Get student profile
// @Description     	Get student profile by user ID.
// @Tags				StudentProfile
// @Produce			json
// @Security 		ApiKeyAuth
// @Success			200			{object}		studentprofile.StudentProfileDetails
// @Failure			400			{object}		utils.ApiResponses	"Invalid request"
// @Failure			401			{object}		utils.ApiResponses	"Unauthorized"
// @Failure			500			{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/student-profile	[get]
func (db Database) ReadStudentProfile(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	profile, err := domains.ReadStudentProfile(db.DB, session.UserID)
	if err != nil {
		// Graceful handling of record not found
		if err.Error() == "record not found" {
			logrus.Warnf("Student profile not found for user_id: %s", session.UserID)
			utils.BuildResponse(ctx, http.StatusOK, "No profile found", nil)
			return
		}

		logrus.WithError(err).Errorf("Unexpected error retrieving student profile for user_id: %s", session.UserID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, profile)
}


// UpdateStudentProfile 	Handles the update of a student profile.
// @Summary         	Update student profile
// @Description     	Update student profile by user ID.
// @Tags				StudentProfile
// @Accept				json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param				request	body			studentprofile.StudentProfileIn	true	"Student profile input"
// @Success			200			{object}		utils.ApiResponses
// @Failure			400			{object}		utils.ApiResponses	"Invalid request"
// @Failure			401			{object}		utils.ApiResponses	"Unauthorized"
// @Failure			500			{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/student-profile	[put]
func (db Database) CreateOrUpdateStudentProfile(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(domains.StudentProfile)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid input. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	input.UserID = session.UserID
	input.Step = true // <-- add this

	// Check if profile exists
	existing, err := domains.ReadStudentProfile(db.DB, session.UserID)
	if err != nil && err.Error() != "record not found" {
		logrus.Error("Error checking profile existence: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if existing != nil && existing.ID != uuid.Nil {
		// Profile exists → update it
		if err := db.DB.Model(&domains.StudentProfile{}).
			Where("user_id = ?", session.UserID).
			Updates(input).Error; err != nil {
			logrus.Error("Failed to update student profile: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Profile updated ",
	Content: "Your profile has been updated ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}
		utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
	} else {
		// Profile doesn't exist → create it
		input.ID = uuid.New()
		if err := db.DB.Create(input).Error; err != nil {
			logrus.Error("Failed to create student profile: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Profile Created",
	Content: "Your profile has been created ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}
		utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
	}


}

// DeleteStudentProfile 	Handles the deletion of a student profile.
// @Summary         	Delete student profile
// @Description     	Delete student profile by user ID.
// @Tags				StudentProfile
// @Produce			json
// @Security 		ApiKeyAuth
// @Success			200			{object}		utils.ApiResponses
// @Failure			400			{object}		utils.ApiResponses	"Invalid request"
// @Failure			401			{object}		utils.ApiResponses	"Unauthorized"
// @Failure			500			{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/student-profile	[delete]
func (db Database) DeleteStudentProfile(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	if err := db.DB.Where("user_id = ?", session.UserID).Delete(&domains.StudentProfile{}).Error; err != nil {
		logrus.Error("Failed to delete student profile. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
// AdminGetStudentProfileByID allows CEO/Manager to fetch a student's profile by user ID.
// @Summary     Admin Get Student Profile by User ID
// @Description Allows CEO or Manager to retrieve a student's profile using user ID.
// @Tags        StudentProfile
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userID path string true "Student User ID"
// @Success     200 {object} studentprofile.StudentProfileDetails
// @Failure     400 {object} utils.ApiResponses
// @Failure     403 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student-profile/{userID} [get]
func (db Database) AdminGetStudentProfileByID(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// ✅ Only CEO or Manager can access
	userRoles, err := domains.ReadUsersRoles(db.DB, session.UserID, session.CompanyID)
	if err != nil || len(userRoles) == 0 {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Role lookup failed", utils.Null())
		return
	}

	isAuthorized := false
	for _, ur := range userRoles {
		roleName, err := domains.ReadRoleName(db.DB, ur.RoleID)
		if err == nil && (roleName == "CEO" || roleName == "Manager") {
			isAuthorized = true
			break
		}
	}
	if !isAuthorized {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Only CEO or Manager can access this.", utils.Null())
		return
	}

	// ✅ Parse user ID from URL
	targetID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// ✅ Fetch student profile
	profile, err := domains.ReadStudentProfile(db.DB, targetID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.BuildErrorResponse(ctx, http.StatusNotFound, "No profile found for this student", utils.Null())
			return
		}
		logrus.WithError(err).Errorf("Error fetching profile for user_id: %s", targetID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, profile)
}
