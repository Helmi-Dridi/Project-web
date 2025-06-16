package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Constants for document types
const (
	DocumentTypePassport         = "passport"
	DocumentTypeTranscript       = "transcript"
	DocumentTypeCV               = "cv"
	DocumentTypeMotivationLetter = "motivation_letter"
	DocumentTypeLanguageTest     = "language_test"
)

// Constants for verification statuses
const (
	VerificationPending  = "pending"
	VerificationApproved = "approved"
	VerificationRejected = "rejected"
)

type StudentDocument struct {
	ID             uuid.UUID  `gorm:"column:id; primaryKey; type:uuid; not null"`
	UserID         uuid.UUID  `gorm:"column:user_id; type:uuid; not null; index"`
	DocumentType   string     `gorm:"type:varchar(50); not null; index"` // passport, cv, etc.
	FilePath       string     `gorm:"type:text; not null"`
	Version        uint       `gorm:"default:1"`
	Submitted      bool       `gorm:"default:false"` // Replaces IsFinalized
	Verification   string     `gorm:"type:varchar(20); default:'pending'"` // pending, approved, rejected
	Remarks        string     `gorm:"type:text"`
	UploadedAt     time.Time  `gorm:"autoCreateTime"`
	IsFinalized    bool      `gorm:"default:false"`                                   // True when submitted by student

	// Optional review tracking
	ReviewedByID   *uuid.UUID `gorm:"type:uuid"` // Admin who reviewed
	ReviewedAt     *time.Time

	gorm.Model
}
// ReadStudentDocuments retrieves all uploaded student documents for a specific user.
func ReadStudentDocuments(db *gorm.DB, userID uuid.UUID) ([]StudentDocument, error) {
	var documents []StudentDocument
	err := db.Where("user_id = ?", userID).Order("document_type, created_at desc").Find(&documents).Error
	return documents, err
}

