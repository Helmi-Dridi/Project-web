package studentdocument

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (db Database) CreateStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	input := new(StudentDocumentIn)

	if err := ctx.ShouldBindJSON(input); err != nil {
		logrus.Error("Invalid input for document upload: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	validDocType := regexp.MustCompile(`^(passport|transcript|cv|motivation_letter|language_test)$`)
	if !validDocType.MatchString(input.DocumentType) || input.FilePath == "" {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc := &domains.StudentDocument{
		ID:           uuid.New(),
		UserID:       session.UserID,
		DocumentType: input.DocumentType,
		FilePath:     input.FilePath,
		Verification: domains.VerificationPending,
	}

	if err := domains.Create(db.DB, doc); err != nil {
		logrus.Error("Database error while creating document: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}

func (db Database) ReadStudentDocuments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page <= 0 || limit <= 0 {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	offset := (page - 1) * limit

	docs, err := ReadLatestPerTypePaginated(db.DB, session.UserID, limit, offset)
	if err != nil {
		logrus.Error("Error retrieving documents: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Count total unique document types for pagination
	var count int64
	err = db.DB.Model(&domains.StudentDocument{}).
		Select("document_type").
		Where("user_id = ?", session.UserID).
		Group("document_type").
		Count(&count).Error
	if err != nil {
		logrus.Error("Error counting document types: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var items []StudentDocumentTable
	for _, d := range docs {
		items = append(items, StudentDocumentTable{
			ID:           d.ID,
			DocumentType: d.DocumentType,
			FilePath:     d.FilePath,
			Verification: d.Verification,
			Remarks:      d.Remarks,
			UploadedAt:   d.UploadedAt,
			Version:      d.Version,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, StudentDocumentPagination{
		Items:      items,
		Page:       uint(page),
		Limit:      uint(limit),
		TotalCount: uint(count),
	})
}


func (db Database) ReadStudentDocumentList(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docs, err := ReadAllList(db.DB, []domains.StudentDocument{}, session.UserID)
	if err != nil {
		logrus.Error("Error retrieving list: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var list []StudentDocumentList
	for _, d := range docs {
		list = append(list, StudentDocumentList{
			ID:           d.ID,
			DocumentType: d.DocumentType,
			FilePath:     d.FilePath,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, list)
}

func (db Database) ReadStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	response := StudentDocumentDetails{
		ID:           doc.ID,
		UserID:       doc.UserID,
		DocumentType: doc.DocumentType,
		FilePath:     doc.FilePath,
		Verification: doc.Verification,
		Remarks:      doc.Remarks,
		UploadedAt:   doc.UploadedAt,
		Version:      doc.Version,
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

func (db Database) DeleteStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	if err := domains.Delete(db.DB, &domains.StudentDocument{}, docID); err != nil {
		logrus.Error("Error deleting: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	
	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "document",
	Content: "Document deleted",
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

func (db Database) UploadStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	documentType := ctx.PostForm("documentType")
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	fileID := uuid.New()
	filename := fileID.String() + "_" + file.Filename
	savePath := "uploads/documents/" + filename
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	doc := &domains.StudentDocument{
		ID:           fileID,
		UserID:       session.UserID,
		DocumentType: documentType,
		FilePath:     savePath,
		Verification: domains.VerificationPending,
	}

	if err := domains.Create(db.DB, doc); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    documentType,
	Content: "Document uploaded ",
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
// ✅ GET /documents/status
func (db Database) GetDocumentStatus(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	documents, err := domains.ReadStudentDocuments(db.DB, session.UserID)
	if err != nil {
		logrus.Error("Failed to fetch documents for status: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	requiredTypes := []string{
		domains.DocumentTypePassport,
		domains.DocumentTypeTranscript,
		domains.DocumentTypeCV,
		domains.DocumentTypeMotivationLetter,
		domains.DocumentTypeLanguageTest,
	}

	docTypeMap := map[string]bool{}
	for _, doc := range documents {
		docTypeMap[doc.DocumentType] = true
	}

	uploadedCount := 0
	for _, docType := range requiredTypes {
		if docTypeMap[docType] {
			uploadedCount++
		}
	}

	progress := float64(uploadedCount) / float64(len(requiredTypes)) * 100
	ctx.JSON(http.StatusOK, gin.H{"progress": progress})
}

// ✅ GET /documents/:id/download
func (db Database) DownloadDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	ctx.File(doc.FilePath)
}

// ✅ PATCH /documents/:id/finalize
func (db Database) FinalizeDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var doc domains.StudentDocument
	doc, err = ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	doc.IsFinalized = true
	db.DB.Save(&doc)
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// ✅ PUT /documents/:id
func (db Database) UpdateStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var input StudentDocumentIn
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	doc.DocumentType = input.DocumentType
	doc.FilePath = input.FilePath
	doc.Verification = domains.VerificationPending // reset verification on update
	db.DB.Save(&doc)
	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    input.DocumentType,
	Content: "Document updated",
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

// ✅ PATCH /admin/documents/:id/verify
func (db Database) AdminVerifyDocument(ctx *gin.Context) {
	docID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	type VerifyIn struct {
		Verification string `json:"verification" binding:"required,oneof=approved rejected"`
		Remarks      string `json:"remarks"`
	}

	var input VerifyIn
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	doc.Verification = input.Verification
	doc.Remarks = input.Remarks
	db.DB.Save(&doc)

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
// GetStudentDocumentStatus returns the completion percentage of required documents.
// @Summary      Document status
// @Description  Returns document upload completion % by document type
// @Tags         Student Documents
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]interface{}
// @Router       /documents/status [get]
func (db Database) GetStudentDocumentStatus(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	requiredDocs := map[string]bool{
		domains.DocumentTypePassport:         false,
		domains.DocumentTypeTranscript:       false,
		domains.DocumentTypeCV:               false,
		domains.DocumentTypeMotivationLetter: false,
		domains.DocumentTypeLanguageTest:     false,
	}

	docs, err := domains.ReadStudentDocuments(db.DB, session.UserID)
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	for _, d := range docs {
		if _, ok := requiredDocs[d.DocumentType]; ok {
			requiredDocs[d.DocumentType] = true
		}
	}

	completed := 0
	for _, ok := range requiredDocs {
		if ok {
			completed++
		}
	}
	percentage := float64(completed) / float64(len(requiredDocs)) * 100

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
		"progress":    percentage,
		"completed":   completed,
		"total":       len(requiredDocs),
		"requirements": requiredDocs,
	})
}
// DownloadStudentDocument allows the user to download their uploaded file.
// @Summary      Download document
// @Description  Download a file by ID
// @Tags         Student Documents
// @Produce      application/octet-stream
// @Security     ApiKeyAuth
// @Param        ID   path  string  true  "Document ID"
// @Success      200
// @Router       /documents/{ID}/download [get]
func (db Database) DownloadStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	doc, err := ReadByID(db.DB, domains.StudentDocument{}, docID)
	if err != nil || doc.UserID != session.UserID {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	
	ctx.FileAttachment(doc.FilePath, "document_"+doc.DocumentType+".pdf")
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "document",
	Content: "Document uploaded ",
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
// FinalizeStudentDocument marks a document as finalized.
// @Summary      Finalize document
// @Description  Mark a student document as ready for review
// @Tags         Student Documents
// @Produce      json
// @Security     ApiKeyAuth
// @Param        ID   path  string  true  "Document ID"
// @Success      200  {object}  utils.ApiResponses
// @Router       /documents/{ID}/finalize [patch]
func (db Database) FinalizeStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	docID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var doc domains.StudentDocument
	if err := db.DB.First(&doc, "id = ? AND user_id = ?", docID, session.UserID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	doc.IsFinalized = true
	if err := db.DB.Save(&doc).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "document",
	Content: "Document Finalized",
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
// VerifyStudentDocument allows CEO or Manager to approve/reject a student document.
// @Summary      Verify document
// @Description  Allows admin to approve or reject a student document
// @Tags         Student Documents
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        ID     path      string               true  "Document ID"
// @Param        body   body      struct {
//               Verification string `json:"verification" binding:"required,oneof=approved rejected"`
//               Remarks      string `json:"remarks"`
//            }               true  "Verification payload"
// @Success      200    {object}  utils.ApiResponses
// @Failure      403    {object}  utils.ApiResponses
// @Failure      400    {object}  utils.ApiResponses
// @Router       /documents/{ID}/verify [patch]
func (db Database) VerifyStudentDocument(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// Check user's roles
	userRoles, err := domains.ReadUsersRoles(db.DB, session.UserID, session.CompanyID)
	if err != nil || len(userRoles) == 0 {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Role lookup failed", utils.Null())
		return
	}

	// Fetch actual role names
	isAuthorized := false
	for _, ur := range userRoles {
		roleName, err := domains.ReadRoleName(db.DB, ur.RoleID)
		if err != nil {
			continue
		}
		if roleName == "CEO" || roleName == "Manager" {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Access denied. Admins only.", utils.Null())
		return
	}

	// Parse document ID
	docID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse input body
	var input struct {
		Verification string `json:"verification" binding:"required,oneof=approved rejected"`
		Remarks      string `json:"remarks"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Fetch document
	var doc domains.StudentDocument
	if err := db.DB.First(&doc, "id = ?", docID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update verification status
	doc.Verification = input.Verification
	doc.Remarks = input.Remarks

	if err := db.DB.Save(&doc).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Document",
	Content: "Document remarks added ",
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
func (db Database) AdminReadStudentDocuments(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	userRoles, err := domains.ReadUsersRoles(db.DB, session.UserID, session.CompanyID)
	if err != nil || len(userRoles) == 0 {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Role lookup failed", utils.Null())
		return
	}

	isAuthorized := false
	for _, ur := range userRoles {
		roleName, err := domains.ReadRoleName(db.DB, ur.RoleID)
		if err != nil {
			continue
		}
		if roleName == "CEO" || roleName == "Manager" {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, "Only CEO or Manager can view this.", utils.Null())
		return
	}

	targetUserID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// ✅ Fetch full document list for admin
	docs, err := ReadAllList(db.DB, []domains.StudentDocument{}, targetUserID)
	if err != nil {
		logrus.Error("Failed to fetch documents: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var result []StudentDocumentTable
	for _, d := range docs {
		result = append(result, StudentDocumentTable{
			ID:           d.ID,
			DocumentType: d.DocumentType,
			FilePath:     d.FilePath,
			IsFinalized:  d.IsFinalized,
			Verification: d.Verification,
			Remarks:      d.Remarks,
			Version:      d.Version,
			UploadedAt:   d.UploadedAt,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, result)
}
