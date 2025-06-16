/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Labels: Represents information about Labels in the system.
		- ID (uuid.UUID): Unique identifier for the Labels.
		- Name (string): Name of the Labels.
		- Description (string): Content of the Labels.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Usage:
	- Import this package to utilize the provided data structures for handling Labels in the project.

	Note:
	- The Labels structure represents information about the ethiquette of each company.

*/

package domains

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Labels struct {
	ID      uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier 
	Name    string    `gorm:"column:type; not null"`                       // name of labels
	Description string    `gorm:"column:content; not null;"` 
	CompanyID	uuid.UUID  `gorm:"column:company_id; type:uuid; not null;"` 
	Color      string       `gorm:"column:type; not null"`            // color of  label
	gorm.Model
}



type ProfileStudent struct {
	ID                     uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`                  // Unique ID
	UserID                 uuid.UUID `gorm:"column:User_id; type:uuid; not null;"`                                    // Linked user
	CompanyID              uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`                                          // Multi-tenancy

	// Personal Info
	DateOfBirth            time.Time `gorm:"column:date_of_birth"`                                             // Optional
	Gender                 string     `gorm:"type:varchar(20);"`                                                // Male, Female, Other
	Nationality            string     `gorm:"type:varchar(100);"`                                               // e.g., Tunisian
	PassportNumber         string     `gorm:"type:varchar(100);"`                                               // Optional
	NationalID             string     `gorm:"type:varchar(100);"`                                               // Optional

	// Contact Details
	PhoneNumber            string     `gorm:"type:varchar(20);not null"`                                        // Required
	AddressLine            string     `gorm:"type:varchar(255);not null"`                                       // Required
	City                   string     `gorm:"type:varchar(100);not null"`                                       // Required
	Country                string     `gorm:"type:varchar(100);not null"`                                       // Required
	ZipCode                string     `gorm:"type:varchar(20);not null"`                                        // Required
	PreferredContactMethod string     `gorm:"type:varchar(20);not null;default:'email'"`                        // 'email', 'whatsapp', 'call'

	// Academic Background
	Qualification          string     `gorm:"type:varchar(100);not null"`                                       // e.g., High School Diploma
	InstitutionName        string     `gorm:"type:varchar(255);not null"`                                       // University / School
	GraduationYear         uint       `gorm:"type:int;default:0"`                                               // Optional
	GPAScore               string     `gorm:"type:varchar(10);"`                                                // GPA or score
	LanguageTestType       string     `gorm:"type:varchar(50);"`                                                // IELTS, TOEFL
	LanguageTestScore      string     `gorm:"type:varchar(10);"`                                                // e.g., 6.5
	CertificateFilePath    string     `gorm:"type:varchar(255);"`                                               // File path or URL to document

	gorm.Model
}


// ReadLabelName reads the name of the Label based on its ID.
func ReadLabelName(db *gorm.DB, LabelID uuid.UUID) (string, error) {
	label := new(Labels)
	err := db.Select("id, name").Where("id = ?", LabelID).First(label).Error
	return label.Name, err
}
