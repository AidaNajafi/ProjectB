package controllers

import (
	"authentication/internal/provider"
	"authentication/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelController struct {
	svc *service.HotelService
}

func NewHotelController(svc *service.HotelService) *HotelController {
	return &HotelController{svc: svc}
}

func (h *HotelController) GetHotelByDateController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		date := ctx.DefaultQuery("date", "")
		if date == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Date parameter is required",
			})
			return
		}
		hotelsByDate, err := h.svc.GetHotelsByDate(ctx, date)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to find hotels in this range of date",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hotel lists based on input date",
			"list":    hotelsByDate,
		})

	}
}

func (h *HotelController) GetHotelByIDController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		inputId := ctx.Param("id")
		id, err := strconv.ParseInt(inputId, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id must be integer!",
			})
			return
		}
		hotelByID, err := h.svc.GetHotelByID(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to get hotel with this id",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hotel by this id:",
			"Hotel":   hotelByID,
		})

	}
}

func (h *HotelController) GetHotelRoomsController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		inputId := ctx.Param("id")
		id, err := strconv.ParseInt(inputId, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id must be integer!",
			})
			return
		}
		date_from := ctx.DefaultQuery("date_from", "")
		date_to := ctx.DefaultQuery("date_to", "")
		if date_from == "" || date_to == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "date_from and date_to are both required",
			})
			return
		}
		hotelRooms, err := h.svc.GetHotelRooms(ctx, id, date_from, date_to)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to get hotel rooms",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hotel list rooms:",
			"rooms":   hotelRooms,
		})
	}
}

type ReservationResponse struct {
	RoomID    int64  `json:"room_id"`
	HotelID   int64  `json:"hotel_id"`
	UserID    int64  `json:"user_id"`
	UserPhone string `json:"user_phone"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	Status    string `json:"status"`
}

type ReservationRequest struct {
	RoomID    int64  `json:"room_id"`
	DateFrom  string `json:"date_from"`
	DateTo    string `json:"date_to"`
	UserID    int64  `json:"user_id"`
	UserPhone string `json:"user_phone"`
}

func (h *HotelController) ReserveHotelController() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body ReservationRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "invalid request body",
			})
			return
		}

		BodySvr := mapHandlerModelToService(body)
		respSvr, err := h.svc.ReserveHotelService(ctx, BodySvr)
		if err != nil {
			var pe *provider.ProviderError
			if errors.As(err, &pe) {
				ctx.JSON(pe.StatusCode, gin.H{"error": string(pe.Body)})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		handlerBody := mapServiceModelToHandlere(respSvr)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Reservation successful",
			"result":  handlerBody,
		})
	}
}

func mapHandlerModelToService(re ReservationRequest) service.ReservationRequest {
	return service.ReservationRequest{
		RoomID:    re.RoomID,
		DateFrom:  re.DateFrom,
		DateTo:    re.DateTo,
		UserID:    re.UserID,
		UserPhone: re.UserPhone,
	}

}

func mapServiceModelToHandlere(res service.ReservationResponse) ReservationResponse {
	return ReservationResponse{
		RoomID:    res.RoomID,
		HotelID:   res.HotelID,
		UserID:    res.UserID,
		UserPhone: res.UserPhone,
		DateFrom:  res.DateFrom,
		DateTo:    res.DateTo,
		Status:    res.Status,
	}
}
