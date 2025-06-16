package studentprofilesettings

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetStudentSettings retrieves the student settings.
// @Summary     Get student settings
// @Description Retrieve student preferences and emergency contact.
// @Tags        StudentSettings
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} StudentSettingsDetails
// @Failure     400 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student/settings [get]
func (db Database) GetStudentSettings(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	settings, err := ReadStudentSettings(db.DB, session.UserID)
	if err != nil {
		// Gracefully handle "record not found"
		if err.Error() == "record not found" {
			logrus.Warnf("Student settings not found for user_id: %s", session.UserID)
			utils.BuildResponse(ctx, http.StatusOK, "No settings found", nil)
			return
		}

		logrus.WithError(err).Errorf("Unexpected error while reading settings for user_id: %s", session.UserID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	response := StudentSettingsDetails{
		ID:                          settings.ID,
		UserID:                      settings.UserID,
		LanguagePreference:          settings.LanguagePreference,
		ReceiveEmailNotifications:   settings.ReceiveEmailNotifications,
		ReceiveSMSNotifications:     settings.ReceiveSMSNotifications,
		AccountDeleted:              settings.AccountDeleted,
		EmergencyName:               settings.EmergencyName,
		EmergencyRelationship:       settings.EmergencyRelationship,
		EmergencyPhoneNumber:        settings.EmergencyPhoneNumber,
		EmergencyEmail:              settings.EmergencyEmail,
		Step: settings.Step,
		CreatedAt:                   settings.CreatedAt,
		UpdatedAt:                   settings.UpdatedAt,
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}


// UpdateStudentSettingsHandler updates student settings.
// @Summary     Update student settings
// @Description Update language preference, notifications, and emergency contact.
// @Tags        StudentSettings
// @Accept      json
// @Produce     json
// @Security    ApiKeyAuth
// @Param       request body StudentSettingsIn true "Updated student settings"
// @Success     200 {object} utils.ApiResponses
// @Failure     400 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student/settings [put]
func (db Database) CreateOrUpdateStudentSettings(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	var input StudentSettingsIn
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Error("Invalid input: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	emailRegex := regexp.MustCompile(constants.EMAIL_REGEX)
	if !emailRegex.MatchString(input.EmergencyEmail) {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Invalid emergency email format", utils.Null())
		return
	}
	input.Step=true
	model := &domains.StudentSettings{
		UserID:                    session.UserID,
		LanguagePreference:        input.LanguagePreference,
		ReceiveEmailNotifications: input.ReceiveEmailNotifications,
		ReceiveSMSNotifications:   input.ReceiveSMSNotifications,
		AccountDeleted:            input.AccountDeleted,
		EmergencyName:             input.EmergencyName,
		EmergencyRelationship:     input.EmergencyRelationship,
		EmergencyPhoneNumber:      input.EmergencyPhoneNumber,
		EmergencyEmail:            input.EmergencyEmail,
		Step: input.Step,
	}

	// 🔁 Upsert logic: Check existence, update if found, else create
	exists, err := CheckStudentSettingsExists(db.DB, session.UserID)
	if err != nil {
		logrus.Error("Error checking settings existence: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if exists {
		if err := UpdateStudentSettings(db.DB, session.UserID, model); err != nil {
			logrus.Error("Failed to update settings: ", err)
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Settings",
	Content: "you updated your Settings card ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}
	} else {
		model.ID = uuid.New()
		if err := CreateStudentSettings(db.DB, model); err != nil {
			logrus.Error("Failed to create settings: ", err)
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
				// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Settings",
	Content: "you created your Settings card ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}


// DeleteStudentSettingsHandler deletes student settings (soft delete).
// @Summary     Delete student settings
// @Description Soft delete student settings
// @Tags        StudentSettings
// @Produce     json
// @Security    ApiKeyAuth
// @Success     200 {object} utils.ApiResponses
// @Failure     400 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student/settings [delete]
func (db Database) DeleteStudentSettingsHandler(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	exists, err := CheckStudentSettingsExists(db.DB, session.UserID)
	if err != nil {
		logrus.Error("Error checking for deletion: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	if !exists {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	if err := db.DB.Where("user_id = ?", session.UserID).Delete(&domains.StudentSettings{}).Error; err != nil {
		logrus.Error("Error deleting: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
// AdminGetStudentSettingsByID allows CEO/Manager to fetch any student's settings.
// @Summary     Admin Get Student Settings by User ID
// @Description Allows admin (CEO or Manager) to retrieve student settings by ID.
// @Tags        StudentSettings
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userID path string true "Student User ID"
// @Success     200 {object} StudentSettingsDetails
// @Failure     400 {object} utils.ApiResponses
// @Failure     403 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student/settings/{userID} [get]
func (db Database) AdminGetStudentSettingsByID(ctx *gin.Context) {
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

	// ✅ Extract and validate target user ID
	targetID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// ✅ Read settings for the given user ID
	settings, err := ReadStudentSettings(db.DB, targetID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.BuildErrorResponse(ctx, http.StatusNotFound, "No settings found for this student", utils.Null())
			return
		}
		logrus.WithError(err).Errorf("Error fetching settings for user_id: %s", targetID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	response := StudentSettingsDetails{
		ID:                          settings.ID,
		UserID:                      settings.UserID,
		LanguagePreference:          settings.LanguagePreference,
		ReceiveEmailNotifications:   settings.ReceiveEmailNotifications,
		ReceiveSMSNotifications:     settings.ReceiveSMSNotifications,
		AccountDeleted:              settings.AccountDeleted,
		EmergencyName:               settings.EmergencyName,
		EmergencyRelationship:       settings.EmergencyRelationship,
		EmergencyPhoneNumber:        settings.EmergencyPhoneNumber,
		EmergencyEmail:              settings.EmergencyEmail,
		CreatedAt:                   settings.CreatedAt,
		UpdatedAt:                   settings.UpdatedAt,
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}
