package universities

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// GetAllUniversities handles GET /universities/{companyID}
// @Summary Get all universities
// @Description Retrieve a paginated list of universities for a company
// @Tags Universities
// @Produce json
// @Security ApiKeyAuth
// @Param companyID path string true "Company ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} universities.UniversityPagination
// @Failure 400 {object} utils.ApiResponses
// @Failure 401 {object} utils.ApiResponses
// @Failure 403 {object} utils.ApiResponses
// @Failure 500 {object} utils.ApiResponses
// @Router /universities/{companyID} [get]
func (db Database) GetAllUniversities(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the page from the request parameter
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(constants.DEFAULT_PAGE_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the limit from the request parameter
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(constants.DEFAULT_LIMIT_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	var universities []domains.University
	if err := db.DB.
		Where("company_id = ?", companyID).
		Offset(offset).
		Limit(limit).
		Find(&universities).Error; err != nil {
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Total count for pagination metadata
	var totalCount int64
	db.DB.Model(&domains.University{}).Where("company_id = ?", companyID).Count(&totalCount)

	var response []UniversityTable
	for _, u := range universities {
		response = append(response, UniversityTable{
			ID:             u.ID,
			Name:           u.Name,
			Country:        u.Country,
			City:           u.City,
			UniversityType: u.UniversityType,
			CreatedAt:      u.CreatedAt,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, UniversityPagination{
		Items:      response,
		Page:       uint(page),
		Limit:      uint(limit),
		TotalCount: uint(totalCount),
	})
}


// GetUniversityPrograms handles GET /universities/{companyID}/{id}/programs
// @Summary Get university programs
// @Description Retrieve all programs for a specific university
// @Tags Universities
// @Produce json
// @Security ApiKeyAuth
// @Param companyID path string true "Company ID"
// @Param id path string true "University ID"
// @Success 200 {array} universities.UniversityProgramInfo
// @Failure 400 {object} utils.ApiResponses
// @Failure 401 {object} utils.ApiResponses
// @Failure 403 {object} utils.ApiResponses
// @Failure 500 {object} utils.ApiResponses
// @Router /universities/{companyID}/{id}/programs [get]
func (db Database) GetUniversityPrograms(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Invalid company UUID: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	universityID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logrus.Error("Invalid university UUID: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Unauthorized access to company data: ", err)
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.UNAUTHORIZED, utils.Null())
		return
	}

	var programs []domains.UniversityProgram
	if err := db.DB.Where("university_id = ?", universityID).Find(&programs).Error; err != nil {
		logrus.Error("Error retrieving university programs: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	var response []UniversityProgramInfo
	for _, p := range programs {
		response = append(response, UniversityProgramInfo{
			ID:                      p.ID,
			ProgramType:             p.ProgramType,
			ProgramName:             p.ProgramName,
			LanguageOfInstruction:   p.LanguageOfInstruction,
			MinimumGPARequirement:   p.MinimumGPARequirement,
			LanguageTestRequirement: p.LanguageTestRequirement,
			OtherTestRequirements:   p.OtherTestRequirements,
			TuitionFeesPerYear:      p.TuitionFeesPerYear,
			EstimatedLivingExpenses: p.EstimatedLivingExpenses,
			VisaRequirements:        p.VisaRequirements,
			ProofOfFinancialMeans:   p.ProofOfFinancialMeans,
			BlockedAccountRequired:  p.BlockedAccountRequired,
			ApplicationDeadline:     p.ApplicationDeadline,
			ApplicationFee:          p.ApplicationFee,
			ScholarshipsAvailable:   p.ScholarshipsAvailable,
			IntakePeriods:           p.IntakePeriods,
			PartTimeWorkAllowed:     p.PartTimeWorkAllowed,
			RecognitionInTunisia:    p.RecognitionInTunisia,
			Notes:                   p.Notes,
			Comments:                p.Comments,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}
