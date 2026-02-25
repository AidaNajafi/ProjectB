package server

import (
	"authentication/internal/controllers"
	"authentication/internal/logger"
	"authentication/internal/service"

	"github.com/gin-gonic/gin"
)

type HotelServer struct {
	HotelCtrl *controllers.HotelController
	jwtSvc    *service.JwtService
	Log       *logger.SlogLogger
}

func NewHotelServer(ctrl *controllers.HotelController, jwt *service.JwtService, log *logger.SlogLogger) *HotelServer {
	return &HotelServer{HotelCtrl: ctrl, jwtSvc: jwt, Log: log}
}

func (h *HotelServer) SetUpHotelRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(controllers.GeneralLog(h.Log))
	protected.Use(controllers.JwtMiddleware(h.jwtSvc))
	hotels := protected.Group("/hotels")
	hotels.GET("/:id", h.HotelCtrl.GetHotelByIDController())
	hotels.GET("/", h.HotelCtrl.GetHotelByDateController())
	hotels.GET("/:id/rooms", h.HotelCtrl.GetHotelRoomsController())
	protected.POST("/reservations", h.HotelCtrl.ReserveHotelController())
}
