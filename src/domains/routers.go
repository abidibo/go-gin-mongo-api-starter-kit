package domains

import (
	"github.com/gin-gonic/gin"
)

// RoutesRegister attaches routes (path + view) to the given gin router group (paths namespace)
func RoutesRegister(router *gin.RouterGroup) {
	router.GET("/", DomainListView)
	router.GET("/:id", DomainDetailView)
	router.POST("/", CreateDomainView)
	router.PUT("/:id", UpdateDomainView)
	router.DELETE("/:id", DeleteDomainView)
}
