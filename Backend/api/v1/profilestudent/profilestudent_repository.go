package profilestudent

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the student profile package.
type Database struct {
	DB *gorm.DB
}

// NewStudentProfileRepository performs automatic migration of student profile-related structures in the database.
func NewStudentProfileRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentProfile{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the student profile structure. Error: ", err)
	}
}

// ReadAllProfilesPagination retrieves a paginated list of student profiles.
func ReadAllProfilesPagination(db *gorm.DB, model []domains.StudentProfile, limit, offset int) ([]domains.StudentProfile, error) {
	err := db.Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllProfilesList retrieves a list of all student profiles.
func ReadAllProfilesList(db *gorm.DB, model []domains.StudentProfile) ([]domains.StudentProfile, error) {
	err := db.Find(&model).Error
	return model, err
}

// ReadStudentProfileByID retrieves a student profile by its unique identifier.
func ReadStudentProfileByID(db *gorm.DB, model domains.StudentProfile, id uuid.UUID) (domains.StudentProfile, error) {
	err := db.First(&model, "id = ?", id).Error
	return model, err
}
