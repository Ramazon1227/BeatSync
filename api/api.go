package api

import (
	docs "github.com/Ramazon1227/BeatSync/api/docs"
	"github.com/Ramazon1227/BeatSync/api/handlers"
	// "github.com/Ramazon1227/BeatSync/api/middleware"
	"github.com/Ramazon1227/BeatSync/config"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           BeatSync API
// @version         1.0
// @description     This is an BeatSync application API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Enter the token with Bearer prefix, e.g. "Bearer abcde12345"

// @security ApiKeyAuth

// SetUpRouter godoc
// @description This is an api gateway
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath        /api/

func SetUpRouter(h handlers.Handler, cfg config.Config) (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	docs.SwaggerInfo.Title = cfg.ServiceName
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.Host = cfg.ServiceHost + cfg.HTTPPort
	docs.SwaggerInfo.Schemes = []string{cfg.HTTPScheme}

	r.Use(customCORSMiddleware())

	api := r.Group("/api")
	api.GET("/ping", h.Ping)
	api.GET("/config", h.GetConfig)
	// v1 := api.Group("/v1")
	{

		// // User
		// v1.POST("user", h.CreateUser)
		// v1.GET("user/:user_id", h.GetUserByID)
		// v1.GET("user", h.GetUserList)
		// v1.PUT("user/:user_id", h.UpdateUser)
		// v1.DELETE("user/:user_id", h.DeleteUser)

		// Auth routes
		// auth := v1.Group("/auth")
		// {
		// 	// auth.POST("/login", h.Login)
		// 	// auth.POST("/logout", h.Logout)
		// }

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

