package server

import (
	"authentication/internal/controllers"
	"authentication/internal/service"

	"github.com/gin-gonic/gin"
)

type HotelServer struct {
	HotelCtrl *controllers.HotelController
	jwtSvc    *service.JwtService
}

func NewHotelServer(ctrl *controllers.HotelController, jwt *service.JwtService) *HotelServer {
	return &HotelServer{HotelCtrl: ctrl, jwtSvc: jwt}
}

func (h *HotelServer) SetUpHotelRoutes(router *gin.Engine) {
	protected := router.Group("/")
	hotels := protected.Group("/hotels")
	hotels.Use(controllers.JwtMiddleware(h.jwtSvc))
	hotels.GET("/:id", h.HotelCtrl.GetHotelByIDController())
	hotels.GET("/", h.HotelCtrl.GetHotelByDateController())
	hotels.GET("/:id/rooms", h.HotelCtrl.GetHotelRoomsController())
	protected.POST("/reservations", h.HotelCtrl.ReserveHotelController())
}
