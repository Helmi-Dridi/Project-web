package universities

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for university repository.
type Database struct {
	DB *gorm.DB
}

// NewUniversityRepository performs automatic migration of University and UniversityProgram structures.
func NewUniversityRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.University{}, &domains.UniversityProgram{}); err != nil {
		logrus.Fatal("Error during auto-migration of university models: ", err)
	}
}

// CreateUniversity adds a new university to the database.
func CreateUniversity(db *gorm.DB, uni *domains.University) error {
	return db.Create(uni).Error
}

// UpdateUniversity updates a university record by ID.
func UpdateUniversity(db *gorm.DB, id uuid.UUID, updated *domains.University) error {
	return db.Model(&domains.University{}).
		Where("id = ?", id).
		Updates(updated).Error
}

// GetUniversityByID retrieves a university by its ID.
func GetUniversityByID(db *gorm.DB, id uuid.UUID) (*domains.University, error) {
	var uni domains.University
	err := db.Preload("Programs").Where("id = ?", id).First(&uni).Error
	return &uni, err
}

// ListUniversities retrieves all universities (optionally can be paginated).
func ListUniversities(db *gorm.DB) ([]domains.University, error) {
	var unis []domains.University
	err := db.Preload("Programs").Find(&unis).Error
	return unis, err
}

// CreateProgram adds a new program under a university.
func CreateProgram(db *gorm.DB, program *domains.UniversityProgram) error {
	return db.Create(program).Error
}

// UpdateProgram updates a program by its ID.
func UpdateProgram(db *gorm.DB, id uuid.UUID, updated *domains.UniversityProgram) error {
	return db.Model(&domains.UniversityProgram{}).
		Where("id = ?", id).
		Updates(updated).Error
}

// GetProgramByID retrieves a specific university program by ID.
func GetProgramByID(db *gorm.DB, id uuid.UUID) (*domains.UniversityProgram, error) {
	var prog domains.UniversityProgram
	err := db.Where("id = ?", id).First(&prog).Error
	return &prog, err
}

// ListProgramsByUniversity retrieves all programs under a specific university.
func ListProgramsByUniversity(db *gorm.DB, universityID uuid.UUID) ([]domains.UniversityProgram, error) {
	var programs []domains.UniversityProgram
	err := db.Where("university_id = ?", universityID).Find(&programs).Error
	return programs, err
}
