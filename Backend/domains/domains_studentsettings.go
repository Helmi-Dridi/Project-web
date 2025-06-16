package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentSettings holds both the student preferences and emergency contact info in one table.
type StudentSettings struct {
	ID                         uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"`  // Unique ID for the settings record
	UserID                     uuid.UUID      `gorm:"column:user_id;type:uuid;not null;uniqueIndex"`                  // User ID this settings belongs to

	// Preferences
	LanguagePreference         string         `gorm:"type:varchar(10);default:'en'"`   // Language code, default is "en"
	ReceiveEmailNotifications  bool           `gorm:"default:true"`                    // Whether student wants email notifications
	ReceiveSMSNotifications    bool           `gorm:"default:false"`                   // Whether student wants SMS notifications
	AccountDeleted             bool           `gorm:"default:false"`                   // Marks if student requested account deletion

	// Emergency Contact
	EmergencyName              string         `gorm:"type:varchar(100);not null"`      // Name of emergency contact
	EmergencyRelationship      string         `gorm:"type:varchar(50);not null"`       // Relationship to student
	EmergencyPhoneNumber       string         `gorm:"type:varchar(20);not null"`       // Emergency phone number
	EmergencyEmail             string         `gorm:"type:varchar(100);not null"`      // Emergency email address
	Step bool `gorm:"column:step;default:false" json:"step"`

	CreatedAt                  time.Time      // Timestamp of record creation
	UpdatedAt                  time.Time      // Timestamp of last update
	DeletedAt                  gorm.DeletedAt `gorm:"index"` // Soft delete support
}
// ReadStudentSettings retrieves the settings and emergency contact info for a specific user.
func ReadStudentSettings(db *gorm.DB, userID uuid.UUID) (*StudentSettings, error) {
	settings := new(StudentSettings)
	err := db.Where("user_id = ?", userID).First(settings).Error
	return settings, err
}
