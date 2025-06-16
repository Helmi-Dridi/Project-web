package studentacademic

import (
	"fmt"
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateStudentAcademic handles the creation of a new student academic background record.
// @Summary         Create academic record
// @Description     Create a new student academic background entry.
// @Tags            StudentAcademic
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @Param           request body studentacademic.StudentAcademicIn true "Academic input"
// @Success         201 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-academic [post]
func (db Database) CreateStudentAcademic(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(domains.StudentAcademicBackground)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid academic input. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	input.ID = uuid.New()
	input.UserID = session.UserID

	if err := db.DB.Create(input).Error; err != nil {
		logrus.Error("Failed to create academic record. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

// ReadStudentAcademic handles the retrieval of a student's academic background by user ID.
// @Summary         Get academic record
// @Description     Retrieve student academic background by user ID.
// @Tags            StudentAcademic
// @Produce         json
// @Security        ApiKeyAuth
// @Success         200 {object} studentacademic.StudentAcademicDetails
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         404 {object} utils.ApiResponses "Not found"
// @Router          /student-academic [get]
func (db Database) ReadStudentAcademic(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	record, err := domains.ReadStudentAcademicBackground(db.DB, session.UserID)
	if err != nil {
		if err.Error() == "record not found" {
			logrus.Warnf("Academic record not found for user_id: %s", session.UserID)
			utils.BuildResponse(ctx, http.StatusOK, "No academic information found", nil)
			return
		}

		logrus.WithError(err).Errorf("Unexpected error retrieving academic record for user_id: %s", session.UserID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, record)
}

// UpdateStudentAcademic handles the update of a student's academic background.
// @Summary         Update academic record
// @Description     Update academic background by user ID.
// @Tags            StudentAcademic
// @Accept          json
// @Produce         json
// @Security        ApiKeyAuth
// @Param           request body studentacademic.StudentAcademicIn true "Academic input"
// @Success         200 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-academic [put]
func (db Database) CreateOrUpdateStudentAcademic(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := &domains.StudentAcademicBackground{}

	// Parse form data (multipart/form-data)
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		logrus.Error("Failed to parse multipart form: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Extract text fields
	input.Qualification = ctx.PostForm("qualification")
	input.InstitutionName = ctx.PostForm("institutionName")
	input.GPAScore = ctx.PostForm("gpaScore")
	input.LanguageTestType = ctx.PostForm("languageTestType")
	input.LanguageTestScore = ctx.PostForm("languageTestScore")
	input.UserID = session.UserID
	input.Step=true
	// Parse GraduationYear safely
	graduationYearStr := ctx.PostForm("graduationYear")
	if yearParsed, err := strconv.ParseUint(graduationYearStr, 10, 32); err == nil {
		input.GraduationYear = uint(yearParsed)
	} else {
		logrus.Error("Invalid graduationYear format: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Handle optional file upload
	file, err := ctx.FormFile("certificate")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		allowed := map[string]bool{
			".pdf":  true,
			".png":  true,
			".jpg":  true,
			".jpeg": true,
			".webp": true,
		}

		if !allowed[ext] {
			logrus.Warn("Unsupported file type: ", ext)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, "Unsupported file type", utils.Null())
			return
		}

		// Save the file
		filePath := fmt.Sprintf("uploads/certificates/%s%s", uuid.New().String(), ext)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			logrus.Error("Failed to save uploaded file: ", err)
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

		// Store path in DB (as relative URL or full URL depending on your frontend)
		input.CertificateFilePath = "/" + filePath
	}

	// Check if record already exists
	existing, err := domains.ReadStudentAcademicBackground(db.DB, session.UserID)
	if err != nil && err.Error() != "record not found" {
		logrus.Error("Error checking academic record existence: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if existing != nil && existing.ID != uuid.Nil {
		// Update
		if err := db.DB.Model(&domains.StudentAcademicBackground{}).
			Where("user_id = ?", session.UserID).
			Updates(input).Error; err != nil {
			logrus.Error("Failed to update academic record. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Academic Background",
	Content: "you updated your academic background ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

		utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
	} else {
		// Create
		input.ID = uuid.New()
		if err := db.DB.Create(input).Error; err != nil {
			logrus.Error("Failed to create academic record. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Academic Background",
	Content: "you created your academic background ",
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



// DeleteStudentAcademic handles deletion of a student's academic background.
// @Summary         Delete academic record
// @Description     Delete academic background by user ID.
// @Tags            StudentAcademic
// @Produce         json
// @Security        ApiKeyAuth
// @Success         200 {object} utils.ApiResponses
// @Failure         400 {object} utils.ApiResponses "Invalid request"
// @Failure         401 {object} utils.ApiResponses "Unauthorized"
// @Failure         500 {object} utils.ApiResponses "Internal Server Error"
// @Router          /student-academic [delete]
func (db Database) DeleteStudentAcademic(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	if err := db.DB.Where("user_id = ?", session.UserID).Delete(&domains.StudentAcademicBackground{}).Error; err != nil {
		logrus.Error("Failed to delete academic record. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Academic Background",
	Content: "you have deleted your academic background ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}


	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
// AdminGetStudentAcademicByID allows CEO/Manager to get a student's academic info by user ID.
// @Summary     Admin Get Academic by User ID
// @Description Allows CEO or Manager to retrieve a student's academic record.
// @Tags        StudentAcademic
// @Produce     json
// @Security    ApiKeyAuth
// @Param       userID path string true "Student User ID"
// @Success     200 {object} studentacademic.StudentAcademicDetails
// @Failure     400 {object} utils.ApiResponses
// @Failure     403 {object} utils.ApiResponses
// @Failure     404 {object} utils.ApiResponses
// @Failure     500 {object} utils.ApiResponses
// @Router      /student-academic/{userID} [get]
func (db Database) AdminGetStudentAcademicByID(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// ✅ Role-based authorization
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
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Only CEO or Manager can access this", utils.Null())
		return
	}

	// ✅ Parse userID from URL
	targetID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// ✅ Fetch academic record
	record, err := domains.ReadStudentAcademicBackground(db.DB, targetID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.BuildErrorResponse(ctx, http.StatusNotFound, "No academic record found", utils.Null())
			return
		}
		logrus.WithError(err).Errorf("Failed to retrieve academic for user_id: %s", targetID)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, record)
}
