package studentcontact

import (
	"time"

	"github.com/google/uuid"
)

// @Description StudentContactIn represents the input structure for creating or updating student contact information.
type StudentContactIn struct {
	PhoneNumber            string `json:"phoneNumber" binding:"required,e164"`                      // Phone number in E.164 format
	Email                  string `json:"email" binding:"required,email"`                           // Student's email address
	AddressLine            string `json:"addressLine"`                                              // Street address
	City                   string `json:"city"`                                                     // City
	Country                string `json:"country"`                                                  // Country
	ZipCode                string `json:"zipCode"`                                                  // ZIP or postal code
	PreferredContactMethod string `json:"preferredContactMethod" binding:"required,oneof=whatsapp email call"` // Preferred contact method
} //@name StudentContactIn

// @Description StudentContactPagination represents a paginated list of student contact details.
type StudentContactPagination struct {
	Items      []StudentContactTable `json:"items"`      // Paginated contact entries
	Page       uint                  `json:"page"`       // Current page
	Limit      uint                  `json:"limit"`      // Page limit
	TotalCount uint                  `json:"totalCount"` // Total record count
} //@name StudentContactPagination

// @Description StudentContactTable represents a student contact entry in a table view.
type StudentContactTable struct {
	ID                     uuid.UUID `json:"id"`                     // Unique contact record ID
	PhoneNumber            string    `json:"phoneNumber"`            // Phone number
	Email                  string    `json:"email"`                  // Email
	City                   string    `json:"city"`                   // City
	Country                string    `json:"country"`                // Country
	PreferredContactMethod string    `json:"preferredContactMethod"` // Preferred communication method
	CreatedAt              time.Time `json:"createdAt"`              // Timestamp of creation
} //@name StudentContactTable

// @Description StudentContactList represents a simplified list view of contact details.
type StudentContactList struct {
	ID          uuid.UUID `json:"id"`          // Contact ID
	PhoneNumber string    `json:"phoneNumber"` // Phone number
	Email       string    `json:"email"`       // Email
} //@name StudentContactList

// @Description StudentContactCount represents the count of contact records.
type StudentContactCount struct {
	Count uint `json:"count"` // Number of contact entries
} //@name StudentContactCount

// @Description StudentContactDetails represents full details of a student contact record.
type StudentContactDetails struct {
	ID                     uuid.UUID `json:"id"`                     // Contact ID
	UserID                 uuid.UUID `json:"userID"`                 // Associated user ID
	PhoneNumber            string    `json:"phoneNumber"`            // Phone number
	Email                  string    `json:"email"`                  // Email
	AddressLine            string    `json:"addressLine"`            // Street address
	City                   string    `json:"city"`                   // City
	Country                string    `json:"country"`                // Country
	ZipCode                string    `json:"zipCode"`                // ZIP code
	PreferredContactMethod string    `json:"preferredContactMethod"` // Preferred communication method
	CreatedAt              time.Time `json:"createdAt"`              // Creation timestamp
	UpdatedAt              time.Time `json:"updatedAt"`              // Last update timestamp
} //@name StudentContactDetails
