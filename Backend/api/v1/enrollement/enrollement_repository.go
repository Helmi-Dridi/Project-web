package studentenrollments

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewStudentEnrollmentRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.StudentEnrollment{}); err != nil {
		logrus.Fatal("Failed to migrate student enrollment model: ", err)
	}
}

func ReadStudentEnrollment(db *gorm.DB, userID uuid.UUID) (*domains.StudentEnrollment, error) {
	enrollment := new(domains.StudentEnrollment)
	err := db.Where("user_id = ?", userID).First(enrollment).Error
	return enrollment, err
}

