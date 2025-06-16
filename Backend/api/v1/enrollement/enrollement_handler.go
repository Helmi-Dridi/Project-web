package studentenrollments

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// AssignStudentToProgram enrolls a student into a specific program.
func (db Database) AssignStudentToProgram(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	userID, _ := uuid.Parse(ctx.Param("userID"))
	programID, _ := uuid.Parse(ctx.Param("programID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	// Ensure user has 'Student' role
	roles, _ := domains.ReadUsersRoles(db.DB, userID, companyID)
	hasStudentRole := false
	for _, r := range roles {
		var role domains.Roles
		db.DB.First(&role, "id = ?", r.RoleID)
		if role.Name == "Student" {
			hasStudentRole = true
			break
		}
	}
	if !hasStudentRole {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, "User is not a student", utils.Null())
		return
	}

	// Validate program
	var program domains.UniversityProgram
	if err := db.DB.First(&program, "id = ?", programID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, "Program not found", utils.Null())
		return
	}

	// Fetch university manually
	var university domains.University
	if err := db.DB.First(&university, "id = ?", program.UniversityID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, "University not found", utils.Null())
		return
	}

	// Remove existing enrollment (if any)
	db.DB.Where("user_id = ?", userID).Delete(&domains.StudentEnrollment{})

	// Create new enrollment
	enrollment := domains.StudentEnrollment{
		ID:        uuid.New(),
		UserID:    userID,
		ProgramID: programID,
	}
	db.DB.Create(&enrollment)

	// Prepare response
	r := EnrollmentResponse{
		UserID:     userID,
		ProgramID:  programID,
		EnrolledAt: enrollment.EnrolledAt,
		Program: ProgramDetails{
			ID:                    program.ID,
			ProgramType:           program.ProgramType,
			ProgramName:           program.ProgramName,
			LanguageOfInstruction: program.LanguageOfInstruction,
		},
		University: UniversityDetails{
			ID:             university.ID,
			Name:           university.Name,
			Country:        university.Country,
			City:           university.City,
			UniversityType: university.UniversityType,
		},
	}
	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Student Assignment :",
	Content: "Student has been assigned to "+university.Name+" in the program "+program.ProgramName,
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, r)
}

// UnassignStudentFromProgram removes a student's enrollment.
func (db Database) UnassignStudentFromProgram(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	userID, _ := uuid.Parse(ctx.Param("userID"))

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	db.DB.Unscoped().Where("user_id = ?", userID).Delete(&domains.StudentEnrollment{})
// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Student Unassignment ",
	Content: "Student has been unassigned to  the program ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, "Student unenrolled successfully")
}

// GetStudentEnrollment retrieves enrollment with program + university details + student info.
func (db Database) GetStudentEnrollment(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// Step 1: Enrollment
	enrollment := new(domains.StudentEnrollment)
	if err := db.DB.Where("user_id = ?", session.UserID).First(enrollment).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, "No enrollment found", utils.Null())
		return
	}

	// Step 2: Program
	program := new(domains.UniversityProgram)
	if err := db.DB.First(program, "id = ?", enrollment.ProgramID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, "Program not found", utils.Null())
		return
	}

	// Step 3: University
	university := new(domains.University)
	if err := db.DB.First(university, "id = ?", program.UniversityID).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusNotFound, "University not found", utils.Null())
		return
	}

	// ✅ Step 4: Student info
	student := new(domains.Users)
	if err := db.DB.Select("first_name", "last_name", "email").Where("id = ?", session.UserID).First(student).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch student info", utils.Null())
		return
	}

	// Step 5: Build Response
	r := EnrollmentResponse{
		UserID:     enrollment.UserID,
		ProgramID:  enrollment.ProgramID,
		EnrolledAt: enrollment.EnrolledAt,
		Program: ProgramDetails{
			ID:                    program.ID,
			ProgramType:           program.ProgramType,
			ProgramName:           program.ProgramName,
			LanguageOfInstruction: program.LanguageOfInstruction,
		},
		University: UniversityDetails{
			ID:             university.ID,
			Name:           university.Name,
			Country:        university.Country,
			City:           university.City,
			UniversityType: university.UniversityType,
		},
		Student: StudentInfo{
			FirstName: student.Firstname,
			LastName:  student.Lastname,
			Email:     student.Email,
		},
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, r)
}

func (db Database) GetEnrollmentsByProgramID(ctx *gin.Context) {
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	programID, _ := uuid.Parse(ctx.Param("programID"))
	session := utils.ExtractJWTValues(ctx)

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	var enrollments []domains.StudentEnrollment
	if err := db.DB.Where("program_id = ?", programID).Find(&enrollments).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch enrollments", utils.Null())
		return
	}

	var responses []EnrollmentResponse
	for _, enrollment := range enrollments {
		var program domains.UniversityProgram
		if err := db.DB.First(&program, "id = ?", enrollment.ProgramID).Error; err != nil {
			continue
		}

		var university domains.University
		if err := db.DB.First(&university, "id = ?", program.UniversityID).Error; err != nil {
			continue
		}

		var student domains.Users
		if err := db.DB.Select("first_name", "last_name", "email").
			Where("id = ?", enrollment.UserID).First(&student).Error; err != nil {
			continue
		}

		responses = append(responses, EnrollmentResponse{
			UserID:     enrollment.UserID,
			ProgramID:  enrollment.ProgramID,
			EnrolledAt: enrollment.EnrolledAt,
			Program: ProgramDetails{
				ID:                    program.ID,
				ProgramType:           program.ProgramType,
				ProgramName:           program.ProgramName,
				LanguageOfInstruction: program.LanguageOfInstruction,
			},
			University: UniversityDetails{
				ID:             university.ID,
				Name:           university.Name,
				Country:        university.Country,
				City:           university.City,
				UniversityType: university.UniversityType,
			},
			Student: StudentInfo{
				FirstName: student.Firstname,
				LastName:  student.Lastname,
				Email:     student.Email,
			},
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, responses)
}

func (db Database) GetAllStudentEnrollments(ctx *gin.Context) {
	companyID, _ := uuid.Parse(ctx.Param("companyID"))
	session := utils.ExtractJWTValues(ctx)

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	var enrollments []domains.StudentEnrollment
	if err := db.DB.Find(&enrollments).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, "Failed to fetch enrollments", utils.Null())
		return
	}

	var responses []EnrollmentResponse
	for _, enrollment := range enrollments {
    var program domains.UniversityProgram
    if err := db.DB.First(&program, "id = ?", enrollment.ProgramID).Error; err != nil {
        continue
    }

    var university domains.University
    if err := db.DB.First(&university, "id = ?", program.UniversityID).Error; err != nil {
        continue
    }

    var user domains.Users
    if err := db.DB.First(&user, "id = ?", enrollment.UserID).Error; err != nil {
        continue
    }

    responses = append(responses, EnrollmentResponse{
        UserID:     enrollment.UserID,
        ProgramID:  enrollment.ProgramID,
        EnrolledAt: enrollment.EnrolledAt,
        Program: ProgramDetails{
            ID:                    program.ID,
            ProgramType:           program.ProgramType,
            ProgramName:           program.ProgramName,
            LanguageOfInstruction: program.LanguageOfInstruction,
        },
        University: UniversityDetails{
            ID:             university.ID,
            Name:           university.Name,
            Country:        university.Country,
            City:           university.City,
            UniversityType: university.UniversityType,
        },
        Student: StudentInfo{
            FirstName: user.Firstname,
            LastName:  user.Lastname,
            Email:     user.Email,
        },
    })
}


	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, responses)
}
