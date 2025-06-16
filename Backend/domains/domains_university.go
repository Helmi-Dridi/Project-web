package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// University represents a university entity.
type University struct {
	ID            uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` 
	Name          string    `gorm:"type:varchar(255);not null"`
	Country       string    `gorm:"type:varchar(100);not null"`
	City          string    `gorm:"type:varchar(100);not null"`
	Ranking       string    `gorm:"type:varchar(100)"`
	UniversityType string   `gorm:"type:varchar(50);not null"` // e.g., "Public" or "Private"
	CompanyID       uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`         // ID of the company to which the user belongs

	Programs      []UniversityProgram `gorm:"foreignKey:UniversityID;constraint:OnDelete:CASCADE"`

	gorm.Model
}


