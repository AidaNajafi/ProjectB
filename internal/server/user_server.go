package server

import (
	"authentication/internal/controllers"

	"github.com/gin-gonic/gin"
)

type UserServer struct {
	UserCtrl *controllers.Controller
}

func NewUserServer(ctrl *controllers.Controller) *UserServer {
	return &UserServer{UserCtrl: ctrl}
}

func (s *UserServer) SetUpUserRoutes(router *gin.Engine) {
	router.POST("/signup", s.UserCtrl.SignUp())
	router.POST("/login", s.UserCtrl.Login())
	router.GET("/health", s.UserCtrl.HealthChecker())

}
