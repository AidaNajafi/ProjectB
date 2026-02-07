package server

import (
	"authentication/internal/controllers"
	"authentication/internal/service"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(authService *service.AuthService) *gin.Engine {
	router := gin.Default()
	router.POST("/signup", signUpHandler(authService))
	router.POST("/login", loginHandler(authService))
	return router
}

func signUpHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		controllers.SignUpGin(c, authService)
	}
}

func loginHandler(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		controllers.LoginGin(c, authService)
	}
}
