package studentprofilesettings

import (
	"time"

	"github.com/google/uuid"
)

// @Description	StudentSettingsIn represents the input structure for saving or updating student settings and emergency contact.
type StudentSettingsIn struct {
	LanguagePreference          string `json:"languagePreference" binding:"required,oneof=en fr de es"` // Language code
	ReceiveEmailNotifications   bool   `json:"receiveEmailNotifications"`                               // Email notification toggle
	ReceiveSMSNotifications     bool   `json:"receiveSMSNotifications"`                                 // SMS notification toggle
	AccountDeleted              bool   `json:"accountDeleted"`                                          // Mark account as deleted
	Step 						bool   `json:"step"`

	EmergencyName               string `json:"emergencyName" binding:"required,min=2,max=100"`          // Emergency contact name
	EmergencyRelationship       string `json:"emergencyRelationship" binding:"required,min=2,max=50"`   // Relationship to the student
	EmergencyPhoneNumber        string `json:"emergencyPhoneNumber" binding:"required,min=6,max=20"`    // Emergency contact phone
	EmergencyEmail              string `json:"emergencyEmail" binding:"required,email"`                 // Emergency contact email
} //@name StudentSettingsIn

// @Description	StudentSettingsDetails represents the detailed structure returned for student settings and emergency contact.
type StudentSettingsDetails struct {
	ID                          uuid.UUID `json:"id"`                          // Settings record ID
	UserID                      uuid.UUID `json:"userID"`                      // Associated user ID

	LanguagePreference          string    `json:"languagePreference"`          // Language code
	ReceiveEmailNotifications   bool      `json:"receiveEmailNotifications"`   // Email notification toggle
	ReceiveSMSNotifications     bool      `json:"receiveSMSNotifications"`     // SMS notification toggle
	AccountDeleted              bool      `json:"accountDeleted"`              // Mark account as deleted
	Step 						bool   `json:"step"`

	EmergencyName               string    `json:"emergencyName"`               // Emergency contact name
	EmergencyRelationship       string    `json:"emergencyRelationship"`       // Relationship to the student
	EmergencyPhoneNumber        string    `json:"emergencyPhoneNumber"`        // Emergency contact phone
	EmergencyEmail              string    `json:"emergencyEmail"`              // Emergency contact email

	CreatedAt                   time.Time `json:"createdAt"`                   // Record creation time
	UpdatedAt                   time.Time `json:"updatedAt"`                   // Last update time
} //@name StudentSettingsDetails
