package universities

import (
	"time"

	"github.com/google/uuid"
)

// @Description UniversityIn represents the input structure for creating or updating a university.
type UniversityIn struct {
	Name           string    `json:"name" binding:"required,min=2,max=255"`           // Name of the university.
	Country        string    `json:"country" binding:"required"`                      // Country where the university is located.
	City           string    `json:"city" binding:"required"`                         // City where the university is located.
	Ranking        string    `json:"ranking"`                                         // Optional university ranking.
	UniversityType string    `json:"universityType" binding:"required"`              // Type: Public or Private.
	CompanyID      uuid.UUID `json:"companyId" binding:"required"`                   // Company that owns this university entry.
} //@name UniversityIn

// @Description UniversityTable is a single entry in a university list view.
type UniversityTable struct {
	ID             uuid.UUID `json:"id"`             // Unique identifier.
	Name           string    `json:"name"`           // Name of the university.
	Country        string    `json:"country"`        // Country of the university.
	City           string    `json:"city"`           // City of the university.
	UniversityType string    `json:"universityType"` // Public or Private.
	CreatedAt      time.Time `json:"createdAt"`      // Date of creation.
} //@name UniversityTable

// @Description UniversityDetails contains all detailed fields for a university including related programs.
type UniversityDetails struct {
	ID             uuid.UUID               `json:"id"`             // ID of the university.
	Name           string                  `json:"name"`           // Name of the university.
	Country        string                  `json:"country"`        // Country where the university is located.
	City           string                  `json:"city"`           // City of the university.
	Ranking        string                  `json:"ranking"`        // University ranking (optional).
	UniversityType string                  `json:"universityType"` // Type: Public or Private.
	Programs       []UniversityProgramInfo `json:"programs"`       // Related study programs.
	CreatedAt      time.Time               `json:"createdAt"`      // Created timestamp.
} //@name UniversityDetails

// @Description UniversityPagination holds paginated list of universities.
type UniversityPagination struct {
	Items      []UniversityTable `json:"items"`      // List of universities.
	Page       uint              `json:"page"`       // Current page number.
	Limit      uint              `json:"limit"`      // Items per page.
	TotalCount uint              `json:"totalCount"` // Total number of universities.
} //@name UniversityPagination

// @Description UniversityProgramIn represents input for creating/updating a program.
type UniversityProgramIn struct {
	
	UniversityID            uuid.UUID `json:"universityId" binding:"required"`         // Related university ID.
	ProgramType             string    `json:"programType" binding:"required"`          // e.g., Bachelor, Master.
	ProgramName             string    `json:"programName" binding:"required"`          // Name of the program.
	LanguageOfInstruction   string    `json:"languageOfInstruction" binding:"required"`// Instruction language.
	MinimumGPARequirement   string    `json:"minimumGpaRequirement"`                   // Optional GPA.
	LanguageTestRequirement string    `json:"languageTestRequirement"`                 // Language test info.
	OtherTestRequirements   string    `json:"otherTestRequirements"`                   // Additional tests.
	TuitionFeesPerYear      string    `json:"tuitionFeesPerYear"`                      // Tuition fees.
	EstimatedLivingExpenses string    `json:"estimatedLivingExpenses"`                 // Living cost estimate.
	VisaRequirements        string    `json:"visaRequirements"`                        // Visa requirement details.
	ProofOfFinancialMeans   string    `json:"proofOfFinancialMeans"`                   // Financial proof.
	BlockedAccountRequired  string    `json:"blockedAccountRequired"`                  // Whether blocked account is needed.
	ApplicationDeadline     string    `json:"applicationDeadline"`                     // Application deadline.
	ApplicationFee          string    `json:"applicationFee"`                          // Fee to apply.
	ScholarshipsAvailable   string    `json:"scholarshipsAvailable"`                   // Scholarship info.
	IntakePeriods           string    `json:"intakePeriods"`                           // Start periods.
	PartTimeWorkAllowed     string    `json:"partTimeWorkAllowed"`                     // Work permissions.
	RecognitionInTunisia    string    `json:"recognitionInTunisia"`                    // If recognized in Tunisia.
	Notes                   string    `json:"notes"`                                   // Additional notes.
	Comments                string    `json:"comments"`                                // Admin/internal comments.
} //@name UniversityProgramIn

// @Description UniversityProgramInfo provides detailed information for a program in responses.
type UniversityProgramInfo struct {
	ID                      uuid.UUID `json:"id"`                      // Unique ID.
	ProgramType             string    `json:"programType"`             // Bachelor, Master, etc.
	ProgramName             string    `json:"programName"`             // Name of the program.
	LanguageOfInstruction   string    `json:"languageOfInstruction"`   // Instruction language.
	MinimumGPARequirement   string    `json:"minimumGpaRequirement"`   // GPA requirement.
	LanguageTestRequirement string    `json:"languageTestRequirement"` // Language test info.
	OtherTestRequirements   string    `json:"otherTestRequirements"`   // Other test info.
	TuitionFeesPerYear      string    `json:"tuitionFeesPerYear"`      // Tuition per year.
	EstimatedLivingExpenses string    `json:"estimatedLivingExpenses"` // Living cost estimate.
	VisaRequirements        string    `json:"visaRequirements"`        // Visa details.
	ProofOfFinancialMeans   string    `json:"proofOfFinancialMeans"`   // Financial proof.
	BlockedAccountRequired  string    `json:"blockedAccountRequired"`  // Blocked account info.
	ApplicationDeadline     string    `json:"applicationDeadline"`     // Deadline to apply.
	ApplicationFee          string    `json:"applicationFee"`          // Application fee.
	ScholarshipsAvailable   string    `json:"scholarshipsAvailable"`   // Scholarships info.
	IntakePeriods           string    `json:"intakePeriods"`           // Intake months/semesters.
	PartTimeWorkAllowed     string    `json:"partTimeWorkAllowed"`     // Work permissions.
	RecognitionInTunisia    string    `json:"recognitionInTunisia"`    // Tunisian recognition.
	Notes                   string    `json:"notes"`                   // Additional notes.
	Comments                string    `json:"comments"`                // Internal comments.
} //@name UniversityProgramInfo
