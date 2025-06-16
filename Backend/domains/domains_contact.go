/*

	Package student_profile provides the data structures for storing and managing student contact information.

	Structures:
	- StudentContactDetails: Contains the contact and communication preferences of a student.
		- ID (uuid.UUID): Unique identifier for the contact record.
		- UserID (uuid.UUID): ID of the user this contact info belongs to.
		- PhoneNumber (string): Student's phone number (E.164 format recommended).
		- Email (string): Student's primary email address.
		- AddressLine (string): Address or street line.
		- City (string): City of residence.
		- Country (string): Country of residence.
		- ZipCode (string): Postal or ZIP code.
		- PreferredContactMethod (PreferredContactMethod): Enum representing the student's preferred method of contact (whatsapp, email, call).
		- gorm.Model: Standard GORM model fields (CreatedAt, UpdatedAt, DeletedAt).

	Constants:
	- PreferredContactMethod: Enum values for contact methods.
		- ContactWhatsApp
		- ContactEmail
		- ContactCall

	Dependencies:
	- "github.com/google/uuid": For handling UUIDs.
	- "gorm.io/gorm": ORM for database modeling.

	Usage:
	- Use this domain to manage student contact information, validate communication preferences, and support student outreach workflows.

	Note:
	- Preferred contact methods are restricted to specific allowed values.

	Last update:
	04/06/2025 14:25

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PreferredContactMethod string

const (
	ContactWhatsApp PreferredContactMethod = "whatsapp"
	ContactEmail    PreferredContactMethod = "email"
	ContactCall     PreferredContactMethod = "call"
)

type StudentContactDetails struct {
	ID                     uuid.UUID              `gorm:"column:id; primaryKey; type:uuid; not null;"`
	UserID                 uuid.UUID              `gorm:"column:user_id;type:uuid;not null;uniqueIndex"`
	PhoneNumber            string                 `gorm:"column:phone_number;type:varchar(20);not null"`
	Email                  string                 `gorm:"column:email;type:varchar(100);not null"`
	AddressLine            string                 `gorm:"column:address_line;type:varchar(255)"`
	City                   string                 `gorm:"column:city;type:varchar(100)"`
	Country                string                 `gorm:"column:country;type:varchar(100)"`
	ZipCode                string                 `gorm:"column:zip_code;type:varchar(20)"`
	PreferredContactMethod PreferredContactMethod `gorm:"column:preferred_contact_method;type:varchar(10);check:preferred_contact_method IN ('whatsapp','email','call')"`
	Step bool `gorm:"column:step;default:false" json:"step"`

	gorm.Model
}
// ReadStudentContactDetails retrieves the student contact details based on the given UserID.
func ReadStudentContactDetails(db *gorm.DB, userID uuid.UUID) (*StudentContactDetails, error) {
	contact := new(StudentContactDetails)
	err := db.Where("user_id = ?", userID).First(contact).Error
	return contact, err
}

