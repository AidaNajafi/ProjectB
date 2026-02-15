package server

import (
	"authentication/internal/controllers"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ctrl *controllers.Controller
}

func NewServer(ctrl *controllers.Controller) *Server {
	return &Server{ctrl: ctrl}
}

func (s *Server) SetUpRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/signup", s.ctrl.SignUp())
	router.POST("/login", s.ctrl.Login())
	return router
}
