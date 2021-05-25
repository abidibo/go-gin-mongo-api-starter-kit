package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"os"
	"systems-management-api/auth"
	_ "systems-management-api/core/logger"
	m "systems-management-api/core/middlewares"
	_ "systems-management-api/docs"
	"systems-management-api/domains"
)

func init() {
	// Read settings
	viper.SetConfigFile(fmt.Sprintf("./%s", os.Getenv("APP_SETTINGS")))
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatal(fmt.Sprintf("Error reading settings file, %s", err))
	} else {
		zap.S().Info("Successfully read settings file ./settings.json")
	}
}

// @title Systems Management API
// @version 0.1.0
// @description This is a REST API used to manage Otto systems
// @termsOfService http://swagger.io/terms/

// @contact.name Otto
// @contact.url https://www.otto.to.it
// @contact.email support@otto.to.it

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @authorizationurl /api/auth/login
// @scope.admin Grants read and write access to administrative information

// @license.name MIT
// @license.url https://mit-license.org/

// @host jeeg.otto.to.it:3000
// @BasePath /api
func main() {
	r := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// test append slash, if I call api/v1/domains without trailing slash answers with a 301 to api/v1/domains/ but
	// without cors headers, so it fails
	r.Use(m.CORSMiddleware())
	r.Use(auth.AuthenticationMiddleware)
	api := r.Group("/api")
	auth.RoutesRegister(api.Group("/auth"))
	domains.RoutesRegister(api.Group("/domain"))
	r.Run()
}
