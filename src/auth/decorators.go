package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"systems-management-api/core/utils"
)

func LoginRequired(view func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		_, exists := c.Get("user")
		if !exists {
			zap.S().Debug("LoginRequired: access to requested view without permission")
			c.JSON(http.StatusForbidden, utils.ErrorResponse{
				Message: "You don't have the rights to see the requested content",
			})
			return
		}
		view(c)
	}
}

func RoleRequired(roles []string, view func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		iuser, exists := c.Get("user")
		if !exists || iuser.(*User).isAnonymous() || !utils.Contains(roles, iuser.(*User).Role) {
			zap.S().Debug("RoleRequired: access to requested view without permission")
			c.JSON(http.StatusForbidden, utils.ErrorResponse{
				Message: "You don't have the rights to see the requested content",
			})
			return
		}
		view(c)
	}
}
