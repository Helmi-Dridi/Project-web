/*

	Package student_profile provides the data structures representing the profile information of students in the system.

	Structures:
	- StudentProfile: Contains personal profile details of the student.
		- ID (uuid.UUID): Unique identifier for the student profile.
		- UserID (uuid.UUID): ID of the user this profile is associated with.
		- DateOfBirth (*time.Time): Optional field for date of birth.
		- Gender (Gender): Enum type representing gender (male, female, other).
		- Nationality (string): The student's country of nationality.
		- PassportNumber (string): Student's passport number (optional).
		- NationalID (string): Student’s national identification number.
		- gorm.Model: Standard GORM model fields (CreatedAt, UpdatedAt, DeletedAt).

	Constants:
	- Gender: Enum values for gender.
		- Male
		- Female
		- Other

	Dependencies:
	- "github.com/google/uuid": Package for handling UUIDs.
	- "gorm.io/gorm": GORM ORM library.
	- "time": Standard Go package for time types.

	Usage:
	- Import this package to manage student profile information in the system.

	Note:
	- Each student profile is tied to a single user via the UserID.
	- Gender is limited to one of the predefined enum values.

	Last update:
	04/06/2025 14:25

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

type StudentProfile struct {
	ID             uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`
	UserID         uuid.UUID `gorm:"column:user_id;type:uuid;not null;uniqueIndex"`
	DateOfBirth    *time.Time `gorm:"column:date_of_birth"`
	Gender         Gender    `gorm:"column:gender;type:varchar(10);check:gender IN ('male','female','other')"`
	Nationality    string    `gorm:"column:nationality;type:varchar(100);not null"`
	PassportNumber string    `gorm:"column:passport_number;type:varchar(50)"`
	NationalID     string    `gorm:"column:national;type:varchar(50);not null"`
	Step bool `gorm:"column:step;default:false" json:"step"`

	gorm.Model
}
// ReadStudentProfile retrieves the student profile based on the given UserID.
func ReadStudentProfile(db *gorm.DB, userID uuid.UUID) (*StudentProfile, error) {
	profile := new(StudentProfile)
	err := db.Where("user_id = ?", userID).First(profile).Error
	return profile, err
}