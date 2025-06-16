package studentprofilesettings

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for student settings.
type Database struct {
	DB *gorm.DB
}

// NewStudentSettingsRepository performs automatic migration of student settings structure.
func NewStudentSettingsRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentSettings{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the student settings structure. Error: ", err)
	}
}

// ReadStudentSettings retrieves the student settings for a specific user.
func ReadStudentSettings(db *gorm.DB, userID uuid.UUID) (*domains.StudentSettings, error) {
	settings := new(domains.StudentSettings)
	err := db.Where("user_id = ?", userID).First(settings).Error
	return settings, err
}

// CreateStudentSettings creates a new student settings record.
func CreateStudentSettings(db *gorm.DB, model *domains.StudentSettings) error {
	return db.Create(model).Error
}

// UpdateStudentSettings updates an existing student settings record by user ID.
func UpdateStudentSettings(db *gorm.DB, userID uuid.UUID, updated *domains.StudentSettings) error {
	return db.Model(&domains.StudentSettings{}).
		Where("user_id = ?", userID).
		Updates(updated).Error
}

// CheckStudentSettingsExists checks whether a settings record exists for the given user.
func CheckStudentSettingsExists(db *gorm.DB, userID uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&domains.StudentSettings{}).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}
