package studentdocument

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database holds the DB instance for use in handlers.
type Database struct {
	DB *gorm.DB
}

// NewStudentDocumentRepository migrates the StudentDocument model.
func NewStudentDocumentRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentDocument{}); err != nil {
		logrus.Fatal("Error migrating StudentDocument structure: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of student documents.
func ReadAllPagination(db *gorm.DB, model []domains.StudentDocument, userID uuid.UUID, limit, offset int) ([]domains.StudentDocument, error) {
	err := db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at desc").Find(&model).Error
	return model, err
}

// ReadAllList retrieves all student documents for a given user.
func ReadAllList(db *gorm.DB, model []domains.StudentDocument, userID uuid.UUID) ([]domains.StudentDocument, error) {
	err := db.Where("user_id = ?", userID).Order("created_at desc").Find(&model).Error
	return model, err
}

// ReadByID retrieves a specific document by its ID.
func ReadByID(db *gorm.DB, model domains.StudentDocument, id uuid.UUID) (domains.StudentDocument, error) {
	err := db.First(&model, id).Error
	return model, err
}

func ReadLatestPerTypePaginated(db *gorm.DB, userID uuid.UUID, limit, offset int) ([]domains.StudentDocument, error) {
	var docs []domains.StudentDocument

	subquery := db.
		Select("MAX(created_at)").
		Table("student_documents AS sub").
		Where("sub.user_id = student_documents.user_id AND sub.document_type = student_documents.document_type")

	err := db.
		Table("student_documents").
		Where("user_id = ?", userID).
		Where("created_at = (?)", subquery).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&docs).Error

	return docs, err
}

