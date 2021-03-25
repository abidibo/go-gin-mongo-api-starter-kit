package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// *ModelValidator containing two parts:
// - Validator: write the json checking rule according to the doc https://github.com/go-playground/validator
// - Model: fill with data from Validator after invoking Bind
// Then, you can just call model.save() after the data is ready in DataModel.
type UserValidatorData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=255"`
	Role     string `json:"role" binding:"required,oneof=superadmin admin"`
}
type UserValidator struct {
	UserData UserValidatorData `json:"user"`
	user     User              `json:"-"`
}

type UserUpdateValidatorData struct {
	Email    string `json:"email,omitempty" binding:"email"`
	Password string `json:"password,omitempty" binding:"min=8,max=255"`
	Role     string `json:"role,omitempty" binding:"oneof=superadmin admin"`
}
type UserUpdateValidator struct {
	UserUpdateData UserUpdateValidatorData `json:"user"`
	user           User                    `json:"-"`
}

// There are some difference when you create or update a model, you need to fill the DataModel before
// update so that you can use your origin data to cheat the validator.
// BTW, you can put your general binding logic here such as setting password.
// @TODO custom error messages https://github.com/gin-gonic/gin/issues/430#issuecomment-141774133
// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go
func (self *UserValidator) Bind(c *gin.Context) error {
	err := c.ShouldBind(&self.UserData)
	if err != nil {
		zap.S().Debug("User Validation Error: ", err)
		return err
	}
	self.user.Email = self.UserData.Email
	self.user.Role = self.UserData.Role
	self.user.SetPassword(self.UserData.Password)
	self.user.Created = time.Now().Unix()

	return nil
}

func (self *UserUpdateValidator) BindUpdate(user *User, c *gin.Context) error {
	err := c.ShouldBind(&self.UserUpdateData)
	if err != nil {
		zap.S().Debug("User Validation Error: ", err)
		return err
	}
	self.user.Email = self.UserUpdateData.Email
	self.user.Role = self.UserUpdateData.Role
	if self.UserUpdateData.Password != "" {
		self.user.SetPassword(self.UserUpdateData.Password)
	}
	self.user.Created = user.Created

	return nil
}

// You can put the default value of a Validator here
func NewUserValidator() UserValidator {
	userValidator := UserValidator{}
	return userValidator
}

func NewUserUpdatelValidator(user *User) UserUpdateValidator {
	userValidator := UserUpdateValidator{}
	userValidator.UserUpdateData.Email = user.Email
	userValidator.UserUpdateData.Role = user.Role
	userValidator.UserUpdateData.Password = user.Password
	userValidator.user.ID = user.ID

	return userValidator
}
