package studentcontact

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the student contact package.
type Database struct {
	DB *gorm.DB
}

// NewStudentContactRepository performs automatic migration of student contact-related structures in the database.
func NewStudentContactRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentContactDetails{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the student contact structure. Error: ", err)
	}
}

// ReadAllContactsPagination retrieves a paginated list of student contacts.
func ReadAllContactsPagination(db *gorm.DB, model []domains.StudentContactDetails, limit, offset int) ([]domains.StudentContactDetails, error) {
	err := db.Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadAllContactsList retrieves a list of all student contact details.
func ReadAllContactsList(db *gorm.DB, model []domains.StudentContactDetails) ([]domains.StudentContactDetails, error) {
	err := db.Find(&model).Error
	return model, err
}

// ReadStudentContactByID retrieves a student contact by its unique identifier.
func ReadStudentContactByID(db *gorm.DB, model domains.StudentContactDetails, id uuid.UUID) (domains.StudentContactDetails, error) {
	err := db.First(&model, "id = ?", id).Error
	return model, err
}

// ReadStudentContactByUserID retrieves a student contact by the associated user ID.
func ReadStudentContactByUserID(db *gorm.DB, userID uuid.UUID) (*domains.StudentContactDetails, error) {
	contact := new(domains.StudentContactDetails)
	err := db.Where("user_id = ?", userID).First(contact).Error
	return contact, err
}
