package studentcontact

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateStudentContact handles the creation of a new student contact.
// @Summary         Create student contact
// @Description     Create a new student contact record.
// @Tags            StudentContact
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @Param           request body studentcontact.StudentContactIn true "Student contact input"
// @Success         201 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-contact [post]
func (db Database) CreateStudentContact(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(domains.StudentContactDetails)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid input for student contact. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	input.ID = uuid.New()
	input.UserID = session.UserID

	if err := db.DB.Create(input).Error; err != nil {
		logrus.Error("Failed to create student contact. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadStudentContact handles retrieval of student contact by user ID.
// @Summary         Get student contact
// @Description     Retrieve student contact details by user ID.
// @Tags            StudentContact
// @Produce         json
// @Security        ApiKeyAuth
// @Success         200 {object} studentcontact.StudentContactDetails
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         404 {object} utils.ApiResponses "Not found"
// @Router          /student-contact [get]
func (db Database) ReadStudentContact(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	contact, err := domains.ReadStudentContactDetails(db.DB, session.UserID)
	if err != nil {
		if err.Error() == "record not found" {
			logrus.Warnf("Student contact not found for user_id: %s", session.UserID)
			utils.BuildResponse(ctx, http.StatusOK, "No contact information found", nil)
			return
		}

		logrus.WithError(err).Errorf("Unexpected error retrieving contact for user_id: %s", session.UserID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, contact)
}

// UpdateStudentContact handles updating the student contact record.
// @Summary         Update student contact
// @Description     Update existing student contact by user ID.
// @Tags            StudentContact
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @Param           request body studentcontact.StudentContactIn true "Student contact input"
// @Success         200 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-contact [put]
func (db Database) CreateOrUpdateStudentContact(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(domains.StudentContactDetails)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid input for student contact. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	input.UserID = session.UserID
	input.Step=true
	// Check if contact already exists
	existing, err := domains.ReadStudentContactDetails(db.DB, session.UserID)
	if err != nil && err.Error() != "record not found" {
		logrus.Error("Error checking contact existence: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if existing != nil && existing.ID != uuid.Nil {
		// Update existing contact
		if err := db.DB.Model(&domains.StudentContactDetails{}).
			Where("user_id = ?", session.UserID).
			Updates(input).Error; err != nil {
			logrus.Error("Failed to update contact. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Contact",
	Content: "you updated your Contact ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}
		utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
	} else {
		// Create new contact
		input.ID = uuid.New()
		if err := db.DB.Create(input).Error; err != nil {
			logrus.Error("Failed to create student contact. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Contact",
	Content: "you Created your Contact ",
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

// DeleteStudentContact handles deletion of a student contact record.
// @Summary         Delete student contact
// @Description     Delete student contact by user ID.
// @Tags            StudentContact
// @Produce         json
// @Security        ApiKeyAuth
// @Success         200 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-contact [delete]
func (db Database) DeleteStudentContact(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	if err := db.DB.Where("user_id = ?", session.UserID).Delete(&domains.StudentContactDetails{}).Error; err != nil {
		logrus.Error("Failed to delete contact. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Contact",
	Content: "you have deleted your Contact ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

// Respond with success
utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())

}
// AdminGetStudentContactByID allows CEO/Manager to fetch a student's contact info by user ID.
// @Summary     Admin Get Student Contact by User ID
// @Description Allows CEO or Manager to retrieve a student's contact using user ID.
// @Tags        StudentContact
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userID path string true "Student User ID"
// @Success     200 {object} studentcontact.StudentContactDetails
// @Failure     400 {object} utils.ApiResponses
// @Failure     403 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student-contact/{userID} [get]
func (db Database) AdminGetStudentContactByID(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// ✅ Ensure user is authorized (CEO or Manager)
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

	// ✅ Extract and validate user ID from path
	targetID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// ✅ Fetch contact information for target user
	contact, err := domains.ReadStudentContactDetails(db.DB, targetID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.BuildErrorResponse(ctx, http.StatusNotFound, "No contact information found for this user", utils.Null())
			return
		}
		logrus.WithError(err).Errorf("Failed to retrieve contact for user_id: %s", targetID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, contact)
}
