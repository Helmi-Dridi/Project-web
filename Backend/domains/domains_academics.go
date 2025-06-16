/*

	Package student_profile provides the data structure for tracking academic background of students.

	Structures:
	- StudentAcademicBackground: Captures educational and language qualification data.
		- ID (uuid.UUID): Unique identifier for the academic background record.
		- UserID (uuid.UUID): ID of the user this academic record is associated with.
		- Qualification (string): Name of the degree or diploma.
		- InstitutionName (string): Name of the academic institution.
		- GraduationYear (uint): Year the student graduated.
		- GPAScore (string): GPA or academic score.
		- LanguageTestType (string): Type of language proficiency test (e.g., IELTS, TOEFL).
		- LanguageTestScore (string): Score obtained in the language test.
		- CertificateFilePath (string): File path or URL to the uploaded certificate document.
		- gorm.Model: GORM standard fields (CreatedAt, UpdatedAt, DeletedAt).

	Dependencies:
	- "github.com/google/uuid": UUID utilities.
	- "gorm.io/gorm": GORM ORM support.

	Usage:
	- Use this domain to maintain student academic records, including GPA and language proficiency details.

	Note:
	- CertificateFilePath may point to local storage or remote URLs depending on configuration.

	Last update:
	04/06/2025 14:25

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentAcademicBackground struct {
	ID                  uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`
	UserID              uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
	Qualification       string    `gorm:"column:qualification;type:varchar(100);not null"`
	InstitutionName     string    `gorm:"column:institution_name;type:varchar(150);not null"`
	GraduationYear      uint      `gorm:"column:graduation_year;check:graduation_year > 1900"`
	GPAScore            string    `gorm:"column:gpa_score;type:varchar(10)"`
	LanguageTestType    string    `gorm:"column:language_test_type;type:varchar(50)"`
	LanguageTestScore   string    `gorm:"column:language_test_score;type:varchar(20)"`
	CertificateFilePath string    `gorm:"column:certificate_file_path;type:text"`
	Step bool `gorm:"column:step;default:false" json:"step"`

	gorm.Model
}
// ReadStudentAcademicBackground retrieves the academic background based on the given UserID.
func ReadStudentAcademicBackground(db *gorm.DB, userID uuid.UUID) (*StudentAcademicBackground, error) {
	academic := new(StudentAcademicBackground)
	err := db.Where("user_id = ?", userID).First(academic).Error
	return academic, err
}

