package studentacademic

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the student academic package.
type Database struct {
	DB *gorm.DB
}

// NewStudentAcademicRepository performs automatic migration of the academic background structure.
func NewStudentAcademicRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentAcademicBackground{}); err != nil {
		logrus.Fatal("An error occurred during migration of the StudentAcademicBackground structure. Error: ", err)
	}
}

// ReadAllAcademicsPagination retrieves a paginated list of academic background records.
func ReadAllAcademicsPagination(db *gorm.DB, model []domains.StudentAcademicBackground, limit, offset int) ([]domains.StudentAcademicBackground, error) {
	err := db.Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllAcademicsList retrieves all academic background records.
func ReadAllAcademicsList(db *gorm.DB, model []domains.StudentAcademicBackground) ([]domains.StudentAcademicBackground, error) {
	err := db.Find(&model).Error
	return model, err
}

// ReadStudentAcademicByID retrieves an academic record by its unique ID.
func ReadStudentAcademicByID(db *gorm.DB, model domains.StudentAcademicBackground, id uuid.UUID) (domains.StudentAcademicBackground, error) {
	err := db.First(&model, "id = ?", id).Error
	return model, err
}

// ReadStudentAcademicByUserID retrieves an academic record by user ID.
func ReadStudentAcademicByUserID(db *gorm.DB, userID uuid.UUID) (*domains.StudentAcademicBackground, error) {
	record := new(domains.StudentAcademicBackground)
	err := db.Where("user_id = ?", userID).First(record).Error
	return record, err
}
