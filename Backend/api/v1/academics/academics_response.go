package studentacademic

import (
	"time"

	"github.com/google/uuid"
)

// @Description StudentAcademicIn represents the input structure for creating or updating student academic records.
type StudentAcademicIn struct {
	Qualification       string `json:"qualification" binding:"required,min=2,max=100"`     // Degree or certificate name
	InstitutionName     string `json:"institutionName" binding:"required,min=2,max=150"`   // Name of the institution
	GraduationYear      uint   `json:"graduationYear" binding:"required,gte=1900"`         // Year of graduation
	GPAScore            string `json:"gpaScore"`                                           // GPA or score
	LanguageTestType    string `json:"languageTestType"`                                   // Type of language test (e.g., IELTS, TOEFL)
	LanguageTestScore   string `json:"languageTestScore"`                                  // Score achieved in the language test
	CertificateFilePath string `json:"certificateFilePath"`                                // File path or URL to the uploaded certificate
} //@name StudentAcademicIn

// @Description StudentAcademicPagination represents the paginated list of student academic records.
type StudentAcademicPagination struct {
	Items      []StudentAcademicTable `json:"items"`      // List of academic entries
	Page       uint                   `json:"page"`       // Current page number
	Limit      uint                   `json:"limit"`      // Max entries per page
	TotalCount uint                   `json:"totalCount"` // Total number of records
} //@name StudentAcademicPagination

// @Description StudentAcademicTable represents a single academic entry in table view.
type StudentAcademicTable struct {
	ID              uuid.UUID `json:"id"`              // Unique academic record ID
	Qualification   string    `json:"qualification"`   // Degree or certificate
	InstitutionName string    `json:"institutionName"` // Institution name
	GraduationYear  uint      `json:"graduationYear"`  // Graduation year
	CreatedAt       time.Time `json:"createdAt"`       // Timestamp of creation
} //@name StudentAcademicTable

// @Description StudentAcademicList represents a simplified list of academic background entries.
type StudentAcademicList struct {
	ID            uuid.UUID `json:"id"`            // Academic record ID
	Qualification string    `json:"qualification"` // Degree or diploma
} //@name StudentAcademicList

// @Description StudentAcademicCount represents the count of academic records.
type StudentAcademicCount struct {
	Count uint `json:"count"` // Total number of academic entries
} //@name StudentAcademicCount

// @Description StudentAcademicDetails represents full details of a student academic record.
type StudentAcademicDetails struct {
	ID                  uuid.UUID `json:"id"`                  // Academic record ID
	UserID              uuid.UUID `json:"userID"`              // Associated user ID
	Qualification       string    `json:"qualification"`       // Degree or diploma
	InstitutionName     string    `json:"institutionName"`     // Name of the institution
	GraduationYear      uint      `json:"graduationYear"`      // Year of graduation
	GPAScore            string    `json:"gpaScore"`            // GPA or academic score
	LanguageTestType    string    `json:"languageTestType"`    // Type of language test
	LanguageTestScore   string    `json:"languageTestScore"`   // Score of the test
	CertificateFilePath string    `json:"certificateFilePath"` // Certificate file location
	CreatedAt           time.Time `json:"createdAt"`           // Creation time
	UpdatedAt           time.Time `json:"updatedAt"`           // Last update time
} //@name StudentAcademicDetails
