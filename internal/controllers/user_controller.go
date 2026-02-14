package controllers

import (
	"authentication/internal/service"
	"log"
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
	jwt *service.JwtService
}

func NewController(svc *service.AuthService, jwt *service.JwtService) *Controller {
	return &Controller{svc: svc, jwt: jwt}
}

var validate = validator.New()

func logf(ctx *gin.Context, format string, args ...any) {
	if v, ok := ctx.Get("logf"); ok {
		if f, ok := v.(func(string, ...any)); ok {
			f(format, args...)
			return
		}
	}
	log.Printf("[authentication]"+format, args...)
}

func (c *Controller) SignUpGin(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body SignUpControl
		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			logf(ctx, "signup bind error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Request",
			})
			return
		}

		if err := validate.Struct(body); err != nil {
			logf(ctx, "field validation error: %v", err)
			ctx.JSON(400, gin.H{
				"message": "All fields must be filled!",
			})
			return
		}
		if len(body.Password) < 8 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "password must be longer than 8 characters!",
			})
			return
		}

		SReq := service.SignUpRequest{
			ID:       body.ID,
			Name:     body.Name,
			Username: body.Username,
			Email:    body.Email,
			Password: body.Password,
		}

		err = c.svc.SignUp(SReq)
		if err != nil {
			logf(ctx, "signup logic failed: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
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
			logf(ctx, "login bind failed: %v", err)
			ctx.JSON(400, gin.H{
				"message": "Invalid request body",
			})
			return
		}
		if err := validate.Struct(body); err != nil {
			logf(ctx, "login field validation failed: %v", err)
			ctx.JSON(400, gin.H{
				"message": "validation failed, username and password are required!",
			})
			return
		}
		SLogReq := service.LoginRequest{
			Username: body.Username,
			Password: body.Password,
		}
		userInfo, err := authService.Login(SLogReq)

		if err != nil {
			logf(ctx, "login logic failed: %v", err)
			ctx.JSON(401, gin.H{
				"message": "wrong username or password",
			})
			return
		}
		jwtClaim := service.UserClaim{
			ID:    userInfo.ID,
			Email: userInfo.Email,
		}
		tokenString, err := c.jwt.GenerateToken(jwtClaim)
		if err != nil {
			logf(ctx, "generate token failed: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "something wrong!",
			})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "login successful",
			"user":    body.Username,
			"token":   tokenString,
		})
	}
}
