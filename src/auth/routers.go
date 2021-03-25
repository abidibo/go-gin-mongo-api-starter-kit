package auth

import (
	"github.com/gin-gonic/gin"
)

// RoutesRegister attaches routes (path + view) to the given gin router group (paths namespace)
func RoutesRegister(router *gin.RouterGroup) {
	router.POST("/login", LoginView)
	router.GET("/user/:id", UserDetailView)
	router.GET("/user", UserListView)
	router.POST("/user", CreateUserView)
	router.PUT("/user/:id", UpdateUserView)
	router.DELETE("/user/:id", DeleteUserView)
}
