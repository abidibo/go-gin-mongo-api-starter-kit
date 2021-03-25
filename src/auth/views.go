package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"systems-management-api/core/utils"
)

// LoginCredentials data type for authentication payload
type LoginCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginSuccessResponse struct {
	Token string `json:"token"`
}

// Allows users to authenticate providing email and password
// @Summary Login user
// @Description Generates and sends a jwt token given user credentials (email and password)
// @Tags auth
// @Accept  json
// @Produce  json
// @Param credentials body LoginCredentials true "Email and password"
// @Success 200 {object} LoginSuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func LoginView(ctx *gin.Context) {
	// returned token
	var token string

	// extract credentials from request
	var credential LoginCredentials
	err := ctx.ShouldBindJSON(&credential)

	// context should have credentials!
	if err != nil {
		zap.S().Debug("Login POST request, error binding POST data: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing email or password"})
	} else {
		zap.S().Debugw("Login POST request with provided credentials", "email", credential.Email)
	}

	// check email and password
	var authService AuthenticationService = NewDatabaseAuthenticationService()
	isUserAuthenticated := authService.Authenticate(credential.Email, credential.Password)
	if isUserAuthenticated {
		// generate token
		var jwtService JWTService = JWTAuthService()
		zap.S().Debugw("Authentication was successful, generating token...")
		token = jwtService.GenerateToken(credential.Email)
		zap.S().Debugw("Token generated")
	}

	if token != "" {
		ctx.JSON(http.StatusOK, LoginSuccessResponse{token})
	} else {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse{Message: "Wrong authentication credentials"})
	}
}

// Returns all users, admin or superadmin roles required
// @Summary Users list
// @Description Retrieves all users
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {array} []UserData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/user [get]
func userListView(c *gin.Context) {
	userService := new(UserService)
	users, err := userService.all()

	if err != nil {
		zap.S().Error("Error while getting all users, Reason: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "Cannot fetch users",
		})
	} else {
		serializer := NewUserSerializer()
		c.JSON(http.StatusOK, serializer.SerializeMany(users))
	}
}

var UserListView = RoleRequired([]string{"admin", "superadmin"}, userListView)

// Returns an user given its id
// @Summary Users detail
// @Description Retrieves one user given its id
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} UserData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/user/{id} [get]
func userDetailView(c *gin.Context) {
	userService := new(UserService)
	user, err := userService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting user, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "User not found"})
	} else {
		serializer := NewUserSerializer()
		c.JSON(http.StatusOK, serializer.Serialize(user))
	}
}

var UserDetailView = RoleRequired([]string{"admin", "superadmin"}, userDetailView)

// Creates an user
// @Summary Create user
// @Description Creates an user
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body UserValidatorData true "User data"
// @Success 201 {object} UserData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/user [post]
func createUserView(c *gin.Context) {
	userValidator := NewUserValidator()
	if err := userValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: err.Error()})
		return
	}

	if _, err := userValidator.user.Save(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: fmt.Sprintf("Cannot insert user: %v", err)})
		return
	}
	serializer := NewUserSerializer()
	c.JSON(http.StatusCreated, serializer.Serialize(&userValidator.user))
}

var CreateUserView = RoleRequired([]string{"admin", "superadmin"}, createUserView)

// Updates an user
// @Summary Update user
// @Description Updates an user
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body UserValidatorData true "User data"
// @Success 200 {object} UserData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/user/{id} [put]
func updateUserView(c *gin.Context) {
	userService := new(UserService)
	user, err := userService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting user, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "User not found"})
		return
	}

	userValidator := NewUserUpdatelValidator(user)
	if err := userValidator.BindUpdate(user, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: err.Error()})
		return
	}

	if _, err := userValidator.user.Save(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: fmt.Sprintf("Cannot update user: %v", err)})
		return
	}
	serializer := NewUserSerializer()
	c.JSON(http.StatusOK, serializer.Serialize(&userValidator.user))
}

var UpdateUserView = RoleRequired([]string{"admin", "superadmin"}, updateUserView)

// Deletes an user
// @Summary Delete user
// @Description Deletes an user
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 204
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /auth/user/{id} [delete]
func deleteUserView(c *gin.Context) {
	userService := new(UserService)
	user, err := userService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting user, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "User not found"})
	} else {
		if _, err := user.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Message: err.Error()})
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

var DeleteUserView = RoleRequired([]string{"admin", "superadmin"}, deleteUserView)
