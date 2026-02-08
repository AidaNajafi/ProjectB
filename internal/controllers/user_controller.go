package controllers

import (
	"authentication/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SignUpControl struct {
	ID       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginControl struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type Controller struct {
	svc *service.AuthService
}

func NewController(svc *service.AuthService) *Controller {
	return &Controller{svc: svc}
}

var validate = validator.New()

func (c *Controller) SignUpGin(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body SignUpControl
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "Invalid Request",
			})
			return
		}
		if err := validate.Struct(body); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err,
				"message": "All fields must be filled!",
			})
			return
		}
		SReq := service.SignUpRequest{
			ID:       body.ID,
			Name:     body.Name,
			Username: body.Username,
			Password: body.Password,
		}

		err = c.svc.SignUp(SReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   err,
				"message": "Sign up failed",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Signup successful",
			"user":    body.Username,
		})
	}
}

func (c *Controller) LoginGin(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body LoginControl
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err,
				"message": "Invalid request body",
			})
			return
		}
		if err := validate.Struct(body); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err,
				"message": "validation failed, username and password are required!",
			})
			return
		}
		SLogReq := service.LoginRequest{
			Username: body.Username,
			Password: body.Password,
		}
		err := authService.Login(SLogReq)
		if err != nil {
			ctx.JSON(401, gin.H{
				"error":   err,
				"message": "wrong username or password",
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "login successful",
			"user":    body.Username,
		})
	}
}
