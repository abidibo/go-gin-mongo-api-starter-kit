package auth

import (
	// "h/src/service"
	// "fmt"
	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

// AuthenticationMiddleware adds the user object to the context if a valid token is provided
// If no token or invalid token is provided it adds an anonymous user (an empty user)
func AuthenticationMiddleware(c *gin.Context) {
	const BEARER_SCHEMA = "Bearer"
	authHeader := c.GetHeader("Authorization")
	zap.S().Debug("AuthenticationMiddleware, reading authorization header: ", authHeader)
	anonymousUser := &User{}

	extractedToken := strings.Split(authHeader, "Bearer ")
	if len(extractedToken) != 2 {
		zap.S().Debug("Incorrect format of authentication token")
		c.Set("user", anonymousUser)
		return
	}

	stringToken := strings.TrimSpace(extractedToken[1])
	zap.S().Debug("AuthenticationMiddleware, extracted token: ", stringToken)
	claim, err := JWTAuthService().ValidateToken(stringToken)
	if err != nil {
		zap.S().Debug("AuthenticationMiddleware, invalid token: ", err)
		c.Set("user", anonymousUser)
		return
	}
	zap.S().Debug("AuthenticationMiddleware, found valid token, user: ", claim.Email)
	userService := new(UserService)
	user, err := userService.GetByEmail(claim.Email)
	if err != nil {
		zap.S().Warn("AuthenticationMiddleware, cannot find user associated to valid token, email: ", claim.Email)
		c.Set("user", anonymousUser)
		return
	}
	zap.S().Debug("AuthenticationMiddleware, added user to context: ", claim.Email)
	c.Set("user", user)
	return
}
