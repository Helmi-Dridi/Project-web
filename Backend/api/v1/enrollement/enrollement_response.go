// FILE: studentenrollments/response.go
package studentenrollments

import (
	"time"

	"github.com/google/uuid"
)

type EnrollmentResponse struct {
	UserID     uuid.UUID `json:"userId"`
	ProgramID  uuid.UUID `json:"programId"`
	EnrolledAt time.Time `json:"enrolledAt"`
	Program    ProgramDetails `json:"program"`
	University UniversityDetails `json:"university"`
	Student    StudentInfo     `json:"student"` // ✅ New field

}

type ProgramDetails struct {
	ID                    uuid.UUID `json:"id"`
	ProgramType           string    `json:"programType"`
	ProgramName           string    `json:"programName"`
	LanguageOfInstruction string    `json:"languageOfInstruction"`
}

type UniversityDetails struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	UniversityType string    `json:"universityType"`
}
type StudentInfo struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
}

