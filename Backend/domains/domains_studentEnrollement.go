package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentEnrollment represents a student's enrollment in a university program.
type StudentEnrollment struct {
	ID         uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;" json:"id"`      // Unique identifier for the enrollment
	UserID     uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`                       // The student being enrolled (only one program per user)
	ProgramID  uuid.UUID `gorm:"type:uuid;not null" json:"program_id"`                           // The program the student is enrolled in
	EnrolledAt time.Time `gorm:"autoCreateTime" json:"enrolled_at"`                              // Timestamp of enrollment

	gorm.Model
}
func ReadStudentEnrollmentWithDetails(db *gorm.DB, userID uuid.UUID) (*StudentEnrollment, *UniversityProgram, *University, error) {
	// 1. Get the enrollment
	enrollment := new(StudentEnrollment)
	if err := db.Where("user_id = ?", userID).First(enrollment).Error; err != nil {
		return nil, nil, nil, err
	}

	// 2. Get the program
	program := new(UniversityProgram)
	if err := db.First(program, "id = ?", enrollment.ProgramID).Error; err != nil {
		return enrollment, nil, nil, err
	}

	// 3. Get the university
	university := new(University)
	if err := db.First(university, "id = ?", program.UniversityID).Error; err != nil {
		return enrollment, program, nil, err
	}

	return enrollment, program, university, nil
}
