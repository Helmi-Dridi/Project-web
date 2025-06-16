package studentdocument

import (
	"time"

	"github.com/google/uuid"
)

// @Description StudentDocumentIn represents the input structure for uploading a document.
type StudentDocumentIn struct {
	DocumentType string `json:"documentType" binding:"required,oneof=passport transcript cv motivation_letter language_test"` // Type of the document
	FilePath     string `json:"filePath" binding:"required"`                                                                  // File path or URL of the uploaded document
	IsFinalized  bool   `json:"isFinalized"`                                                                                  // Indicates whether the document is finalized
} //@name StudentDocumentIn

// @Description StudentDocumentReviewInput is used by admins to review a document.
type StudentDocumentReviewInput struct {
	Verification string `json:"verification" binding:"required,oneof=approved rejected"` // Admin decision
	Remarks      string `json:"remarks"`                                                 // Optional feedback or reason
} //@name StudentDocumentReviewInput

// @Description StudentDocumentTable is a minimal row for listing uploaded documents.
type StudentDocumentTable struct {
	ID           uuid.UUID `json:"id"`
	DocumentType string    `json:"documentType"`
	FilePath     string    `json:"filePath"`
	IsFinalized  bool      `json:"isFinalized"`
	Verification string    `json:"verification"` // pending, approved, rejected
	Remarks      string    `json:"remarks"`
	Version      uint      `json:"version"`
	UploadedAt   time.Time `json:"uploadedAt"`
} //@name StudentDocumentTable

// @Description StudentDocumentPagination returns paginated document list.
type StudentDocumentPagination struct {
	Items      []StudentDocumentTable `json:"items"`
	Page       uint                   `json:"page"`
	Limit      uint                   `json:"limit"`
	TotalCount uint                   `json:"totalCount"`
} //@name StudentDocumentPagination

// @Description StudentDocumentDetails gives full document info for detailed views.
type StudentDocumentDetails struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userID"`
	DocumentType string    `json:"documentType"`
	FilePath     string    `json:"filePath"`
	IsFinalized  bool      `json:"isFinalized"`
	Verification string    `json:"verification"`
	Remarks      string    `json:"remarks"`
	Version      uint      `json:"version"`
	UploadedAt   time.Time `json:"uploadedAt"`
} //@name StudentDocumentDetails

// @Description StudentDocumentList returns simplified list of document references.
type StudentDocumentList struct {
	ID           uuid.UUID `json:"id"`
	DocumentType string    `json:"documentType"`
	FilePath     string    `json:"filePath"`
} //@name StudentDocumentList

// @Description StudentDocumentCount returns count of uploaded documents.
type StudentDocumentCount struct {
	Count uint `json:"count"`
} //@name StudentDocumentCount
