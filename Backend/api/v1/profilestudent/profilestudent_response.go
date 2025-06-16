package profilestudent

import (
	"time"

	"github.com/google/uuid"
)

// @Description StudentProfileIn represents the input structure for creating or updating a student profile.
type StudentProfileIn struct {
	DateOfBirth    *time.Time `json:"dateOfBirth"`                                              // Date of birth (optional)
	Gender         string     `json:"gender" binding:"required,oneof=male female other"`        // Gender: male, female, or other
	Nationality    string     `json:"nationality" binding:"required"`                           // Nationality of the student
	PassportNumber string     `json:"passportNumber"`                                           // Passport number (optional)
	NationalID     string     `json:"nationalID" binding:"required"`                            // National identification number
} //@name StudentProfileIn

// @Description StudentProfilePagination represents the paginated list of student profiles.
type StudentProfilePagination struct {
	Items      []StudentProfileTable `json:"items"`      // List of student profile entries
	Page       uint                  `json:"page"`       // Current page number
	Limit      uint                  `json:"limit"`      // Max number of items per page
	TotalCount uint                  `json:"totalCount"` // Total number of student profiles
} //@name StudentProfilePagination

// @Description StudentProfileTable represents a single student profile entry in a table view.
type StudentProfileTable struct {
	ID             uuid.UUID  `json:"id"`             // Unique identifier of the student profile
	Gender         string     `json:"gender"`         // Gender
	Nationality    string     `json:"nationality"`    // Nationality
	NationalID     string     `json:"nationalID"`     // National identification number
	CreatedAt      time.Time  `json:"createdAt"`      // Timestamp of creation
} //@name StudentProfileTable

// @Description StudentProfileList represents a simplified version of student profiles for listing.
type StudentProfileList struct {
	ID          uuid.UUID `json:"id"`          // Unique identifier of the student profile
	Gender      string    `json:"gender"`      // Gender
	Nationality string    `json:"nationality"` // Nationality
} //@name StudentProfileList

// @Description StudentProfileCount represents the count of student profiles.
type StudentProfileCount struct {
	Count uint `json:"count"` // Total number of student profiles
} //@name StudentProfileCount

// @Description StudentProfileDetails represents full details of a student profile.
type StudentProfileDetails struct {
	ID             uuid.UUID  `json:"id"`             // Unique identifier of the student profile
	UserID         uuid.UUID  `json:"userID"`         // Associated user ID
	DateOfBirth    *time.Time `json:"dateOfBirth"`    // Date of birth
	Gender         string     `json:"gender"`         // Gender
	Nationality    string     `json:"nationality"`    // Nationality
	PassportNumber string     `json:"passportNumber"` // Passport number
	NationalID     string     `json:"nationalID"`     // National ID
	CreatedAt      time.Time  `json:"createdAt"`      // Timestamp of creation
} //@name StudentProfileDetails