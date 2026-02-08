package server

import (
	"authentication/internal/controllers"
	"authentication/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *controllers.Controller
}

func NewServer(srv *controllers.Controller) *Server {
	return &Server{srv: srv}
}

func (s *Server) SetUpRoutes(authService *service.AuthService) *gin.Engine {
	router := gin.Default()
	router.POST("/signup", s.srv.SignUpGin(authService))
	router.POST("/login", s.srv.LoginGin(authService))
	return router
}
