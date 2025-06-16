package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)
// UniversityProgram represents a study program within a university.
type UniversityProgram struct {
	ID                        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` 
	UniversityID              uuid.UUID `gorm:"foreignKey:university_id; references:id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProgramType               string    `gorm:"type:varchar(50);not null"`   // e.g., Bachelor, Master
	ProgramName               string    `gorm:"type:varchar(255);not null"`
	LanguageOfInstruction     string    `gorm:"type:varchar(100);not null"`
	MinimumGPARequirement     string    `gorm:"type:varchar(100)"`
	LanguageTestRequirement   string    `gorm:"type:text"`
	OtherTestRequirements     string    `gorm:"type:text"`
	TuitionFeesPerYear        string    `gorm:"type:varchar(100)"`
	EstimatedLivingExpenses   string    `gorm:"type:varchar(100)"`
	VisaRequirements          string    `gorm:"type:text"`
	ProofOfFinancialMeans     string    `gorm:"type:text"`
	BlockedAccountRequired    string    `gorm:"type:varchar(100)"`
	ApplicationDeadline       string    `gorm:"type:varchar(100)"`
	ApplicationFee            string    `gorm:"type:varchar(100)"`
	ScholarshipsAvailable     string    `gorm:"type:text"`
	IntakePeriods             string    `gorm:"type:varchar(100)"`
	PartTimeWorkAllowed       string    `gorm:"type:varchar(100)"`
	RecognitionInTunisia      string    `gorm:"type:varchar(100)"`
	Notes                     string    `gorm:"type:text"`
	Comments                  string    `gorm:"type:text"`

	gorm.Model
}
