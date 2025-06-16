package users

import (
	"fmt"
	"labs/constants"
	"labs/domains"
	"net/http"
	"path/filepath"
	"strconv"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser 		Handles the creation of a new user.
// @Summary        	Create user
// @Description    	Create a new user.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID		path			string				true		"Company ID"
// @Param			request			body			users.UsersIn		true		"User query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses	"Invalid request"
// @Failure			401				{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses	"Forbidden"
// @Failure			500				{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/users/{companyID}	[post]
func (db Database) CreateUser(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a UserIn struct
	user := new(UsersIn)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Hash the user's password
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Create a new user in the database
	dbUser := &domains.Users{
		ID:              uuid.New(),
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		Email:           user.Email,
		Password:        string(hash),
		Status:          true,
		CompanyID:       user.CompanyID,
		CreatedByUserID: session.UserID,
		ProfilePicture: "uploads/user_images/default_img.png",
	}
	if err := domains.Create(db.DB, dbUser); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "user created.",
	Content: "You created a User with a name : "+user.Firstname+" "+user.Lastname+ "now you need to assign a role for it ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

// Respond with success
utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())

}

// ReadUsers 		Handles the retrieval of all users.
// @Summary        	Get users
// @Description    	Get all users.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			page			query		int			false		"Page"
// @Param			limit			query		int			false		"Limit"
// @Param			companyID		path		string		true		"Company ID"
// @Success			200				{object}	users.UsersPagination
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/users/{companyID}	[get]
func (db Database) ReadUsers(ctx *gin.Context) {

	// Extract JWT values from the context
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

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the user's value is among the allowed choices
	validChoices := utils.ResponseLimitPagination()
	isValidChoice := false
	for _, choice := range validChoices {
		if uint(limit) == choice {
			isValidChoice = true
			break
		}
	}

	// If the value is invalid, set it to default DEFAULT_LIMIT_PAGINATION
	if !isValidChoice {
		limit = constants.DEFAULT_LIMIT_PAGINATION
	}

	// Generate offset
	offset := (page - 1) * limit

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all user data from the database
	users, err := ReadAllPagination(db.DB, []domains.Users{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Users{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a user structure as a response
	response := UsersPagination{}
	dataTableUser := []UsersTable{}
	for _, user := range users {

		dataTableUser = append(dataTableUser, UsersTable{
			ID:        user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
	response.Items = dataTableUser
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadUsersList 	Handles the retrieval the list of all users.
// @Summary        	Get list of  users
// @Description    	Get list of all users.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Success			200					{array}			users.UsersList
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/users/{companyID}/list	[get]
func (db Database) ReadUsersList(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all user data from the database
	users, err := ReadAllList(db.DB, []domains.Users{}, session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a user structure as a response
	usersList := []UsersList{}
	for _, user := range users {
		usersList = append(usersList, UsersList{
			ID:   user.ID,
			Name: user.Firstname + " " + user.Lastname,
		})
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, usersList)
}

// ReadUsersCount 	Handles the retrieval the number of all users.
// @Summary        	Get number of  users
// @Description    	Get number of all users.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID				path			string		true	"Company ID"
// @Success			200						{object}		users.UsersCount
// @Failure			400						{object}		utils.ApiResponses	"Invalid request"
// @Failure			401						{object}		utils.ApiResponses	"Unauthorized"
// @Failure			403						{object}		utils.ApiResponses	"Forbidden"
// @Failure			500						{object}		utils.ApiResponses	"Internal Server Error"
// @Router			/users/{companyID}/count	[get]
func (db Database) ReadUsersCount(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve all user data from the database
	users, err := domains.ReadTotalCount(db.DB, &[]domains.Users{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding all user data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a user structure as a response
	usersCount := UsersCount{
		Count: users,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, usersCount)
}

// ReadUser 		Handles the retrieval of one user.
// @Summary        	Get user
// @Description    	Get one user.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"User ID"
// @Success			200					{object}		users.UsersDetails
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/users/{companyID}/{ID}	[get]
func (db Database) ReadUser(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve user data from the database
	user, err := ReadByID(db.DB, domains.Users{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving user data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a user structure as a response
	details := UsersDetails{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Country:   user.Country,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateUser 		Handles the update of a user.
// @Summary        	Update user
// @Description    	Update user.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string				true	"Company ID"
// @Param			ID					path			string				true	"User ID"
// @Param			request				body			users.UsersIn		true	"User query params"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses			"Invalid request"
// @Failure			401					{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses			"Forbidden"
// @Failure			500					{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/users/{companyID}/{ID}	[put]
func (db Database) UpdateUser(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a UserIn struct
	user := new(UsersIn)
	if err := ctx.ShouldBindJSON(user); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the user with the specified ID exists
	if err = domains.CheckByID(db.DB, &domains.Users{}, objectID); err != nil {
		logrus.Error("Error checking if the user with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Update the user data in the database
	dbUser := &domains.Users{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
	}
	if err = domains.Update(db.DB, dbUser, objectID); err != nil {
		logrus.Error("Error updating user data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Update",
	Content: "You created a User now you need to assign a role for it ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

// Respond with success
utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
	
}

// DeleteUser	 	Handles the deletion of a user.
// @Summary        	Delete user
// @Description    	Delete one user.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"User ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/users/{companyID}/{ID}	[delete]
func (db Database) DeleteUser(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the user with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Users{}, objectID); err != nil {
		logrus.Error("Error checking if the user with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the user data from the database
	if err := domains.Delete(db.DB, &domains.Users{}, objectID); err != nil {
		logrus.Error("Error deleting user data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
		// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Delete",
	Content: "You have deleted a User ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

// Respond with success
utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())

}
// AssignRole 		Handles the assignment of a role to a user.
// @Summary        	Assign role
// @Description    	Assign a role to a user.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Param			ID					path			string			true	"User ID"
// @Param			roleID				path			string			true	"Role ID"
// @Success			200					{object}		utils.ApiResponses
// @Failure			400					{object}		utils.ApiResponses		"Invalid request"
// @Failure			401					{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403					{object}		utils.ApiResponses		"Forbidden"
// @Failure			500					{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/users/{companyID}/{ID}/roles/{roleID}	[post]
func (db Database) AssignRole(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	roleID, err := uuid.Parse(ctx.Param("roleID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Check if the user has the Manager role
	hasManagerRole := false
	for _, role := range session.Roles {
		if role.Name == "CEO" && role.CompanyID == session.CompanyID {
			hasManagerRole = true
			break
		}
	}

	if !hasManagerRole {
		logrus.Error("Unauthorized: Only Managers can perform this action")
		utils.BuildErrorResponse(ctx, http.StatusForbidden, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Check if the user with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Users{}, objectID); err != nil {
		logrus.Error("Error checking if the user with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	//check if the role with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Roles{}, roleID); err != nil {
		logrus.Error("Error checking if the role with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Assign the role to the user
	dbUserRoles := &domains.UsersRoles{
		UserID:    objectID,
		RoleID:    roleID,
		CompanyID: companyID,
		
	}
	if err := domains.Create(db.DB, dbUserRoles); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Assign Role",
	Content: "Role has been assigned to a User ",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}

// Respond with success
utils.BuildResponse(ctx, http.StatusCreated, constants.CREATED, utils.Null())
}


// UploadProfileImage handles uploading and saving a user's profile image.
// @Summary         Upload user profile image
// @Description     Upload a new profile image for the user.
// @Tags            Users
// @Accept          multipart/form-data
// @Produce         json
// @Security        ApiKeyAuth
// @Param           companyID  path      string              true   "Company ID"
// @Param           ID         path      string              true   "User ID"
// @Param           file       formData  file                true   "Profile Picture"
// @Success         200        {object}  utils.ApiResponses
// @Failure         400        {object}  utils.ApiResponses  "Invalid request"
// @Failure         401        {object}  utils.ApiResponses  "Unauthorized"
// @Failure         403        {object}  utils.ApiResponses  "Forbidden"
// @Failure         500        {object}  utils.ApiResponses  "Internal Server Error"
// @Router          /users/{companyID}/{ID}/upload-profile-image [post]
func (db Database) UploadProfileImage(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the user ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the uploaded file
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		logrus.Error("Error retrieving file from request: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	
	// Define the destination directory and file path
	destDir := "./uploads/user_images"
	if err := domains.EnsureDir(destDir); err != nil {
		logrus.Error("Error creating directory: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Set the file name with the user ID to avoid conflicts
	fileName := fmt.Sprintf("%s_%s", objectID, filepath.Base(fileHeader.Filename))
	destPath := filepath.Join(destDir, fileName)

	// Save the uploaded image
	if _, err :=domains.SaveImageFile(fileHeader, destPath); err != nil {
		logrus.Error("Error saving profile picture: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Update the user's profile picture path in the database
	if err := UpdateProfilePicture(db.DB, objectID, destPath); err != nil {
		logrus.Error("Error updating user profile picture in database: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	
			// ✅ Create a notification for the new user
notif := &domains.Notifications{
	ID:      uuid.New(),
	Type:    "Image upload",
	Content: "An Image has been uploaded",
	Seen:    false,
	UserID:  session.UserID,
}
if err := domains.Create(db.DB, notif); err != nil {
	logrus.Error("Failed to create notification for user: ", err)
	// Optional: don’t fail the whole request, just log
}


	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{"profile_picture": destPath})
}
// ReadAdminsList handles the retrieval of all users with the role CEO or Manager.
// @Summary         Get Admins List
// @Description     Get users who have 'CEO' or 'Manager' roles
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Success			200	{array}		users.UsersList
// @Failure			400	{object}	utils.ApiResponses
// @Failure			401	{object}	utils.ApiResponses
// @Failure			403	{object}	utils.ApiResponses
// @Failure			500	{object}	utils.ApiResponses
// @Router			/admins/list [get]
func (db Database) ReadAdminsList(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Invalid UUID format: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Employee does not belong to company: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var users []domains.Users
	if err := db.DB.
		Joins("JOIN users_roles ur ON ur.users_id = users.id").
		Joins("JOIN roles r ON r.id = ur.roles_id").
		Where("r.name IN ?", []string{"CEO", "Manager"}).
		Where("users.status = ?", true).
		Where("users.company_id = ?", companyID).
		Find(&users).Error; err != nil {
		logrus.Error("Error retrieving admin users: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	adminList := make([]UsersList, 0, len(users))
	for _, user := range users {
		adminList = append(adminList, UsersList{
			ID:   user.ID,
			Name: user.Firstname + " " + user.Lastname,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, adminList)
}
// ReadAdminsList handles the retrieval of all users with the role CEO or Manager.
// @Summary         Get Admins List
// @Description     Get users who have 'CEO' or 'Manager' roles
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID			path			string			true	"Company ID"
// @Success			200	{array}		users.AdminList
// @Failure			400	{object}	utils.ApiResponses
// @Failure			401	{object}	utils.ApiResponses
// @Failure			403	{object}	utils.ApiResponses
// @Failure			500	{object}	utils.ApiResponses
// @Router			/admins/list/{companyID} [get]
func (db Database) ReadAdminsListall(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Invalid UUID format: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Employee does not belong to company: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var users []domains.Users
	if err := db.DB.
		Joins("JOIN users_roles ur ON ur.users_id = users.id").
		Joins("JOIN roles r ON r.id = ur.roles_id").
		Where("r.name IN ?", []string{"CEO", "Manager"}).
		Where("users.status = ?", true).
		Where("users.company_id = ?", companyID).
		Find(&users).Error; err != nil {
		logrus.Error("Error retrieving admin users: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	adminList := make([]AdminList, 0, len(users))
	for _, user := range users {
		adminList = append(adminList, AdminList{
			ID: user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Country:   user.Country,
			Status:    user.Status,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, adminList)
}
// ReadUsersWithNoRole handles retrieval of users with no roles.
// @Summary         Get Users Without Roles
// @Description     Get all users in a company who do not have any roles assigned.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	path	string	true	"Company ID"
// @Success			200	{array}	users.AdminList
// @Failure			400	{object}	utils.ApiResponses
// @Failure			401	{object}	utils.ApiResponses
// @Failure			403	{object}	utils.ApiResponses
// @Failure			500	{object}	utils.ApiResponses
// @Router			/users/norole/{companyID} [get]
func (db Database) ReadUsersWithNoRole(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Invalid UUID format: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("User not part of company: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var users []domains.Users
	if err := db.DB.
		Raw(`
			SELECT * FROM users
			WHERE company_id = ? 
			AND status = true 
			AND id NOT IN (SELECT users_id FROM users_roles)
		`, companyID).
		Scan(&users).Error; err != nil {
		logrus.Error("Error retrieving users with no roles: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	result := make([]AdminList, 0, len(users))
	for _, user := range users {
		result = append(result, AdminList{
			ID: user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Country:   user.Country,
			Status:    user.Status,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, result)
}
// ReadStudentsList handles retrieval of users with the 'Student' role.
// @Summary         Get Students List
// @Description     Get all users in a company who have the 'Student' role assigned.
// @Tags			Users
// @Produce			json
// @Security 		ApiKeyAuth
// @Param			companyID	path	string	true	"Company ID"
// @Success			200	{array}	users.AdminList
// @Failure			400	{object}	utils.ApiResponses
// @Failure			401	{object}	utils.ApiResponses
// @Failure			403	{object}	utils.ApiResponses
// @Failure			500	{object}	utils.ApiResponses
// @Router			/users/students/{companyID} [get]
func (db Database) ReadStudentsList(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Invalid UUID format: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("User not part of company: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	var users []domains.Users
	if err := db.DB.
		Joins("JOIN users_roles ur ON ur.users_id = users.id").
		Joins("JOIN roles r ON r.id = ur.roles_id").
		Where("r.name = ?", []string{"Student"}).
		Where("users.status = ?", true).
		Where("users.company_id = ?", companyID).
		Find(&users).Error; err != nil {
		logrus.Error("Error retrieving students: ", err)
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	result := make([]AdminList, 0, len(users))
	for _, user := range users {
		result = append(result, AdminList{
			ID: user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Country:   user.Country,
			Status:    user.Status,
		})
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, result)
}
