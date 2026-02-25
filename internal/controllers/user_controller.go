package controllers

import (
	"authentication/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	svc      *service.AuthService
	validate *validator.Validate
}

func NewUserController(svc *service.AuthService, validator *validator.Validate) *Controller {
	return &Controller{
		svc:      svc,
		validate: validator,
	}
}

type SignUpControl struct {
	ID        int64    `json:"id"`
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	UserPhone string `json:"userphone" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (c *Controller) SignUp() gin.HandlerFunc {
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
		var errMsg string
		if err := c.validate.Struct(body); err != nil {
			logf(ctx, "field validation error: %v", err)
			if ValidationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range ValidationErrors {
					errMsg = getValidationErrorMessage(e)
				}
			}
			ctx.JSON(400, gin.H{
				"error": errMsg,
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
			ID:        body.ID,
			Name:      body.Name,
			Username:  body.Username,
			Email:     body.Email,
			UserPhone: body.UserPhone,
			Password:  body.Password,
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

type LoginControl struct {
	Username  string `json:"username" validate:"required"`
	UserPhone string `json:"userphone validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (c *Controller) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body LoginControl
		if err := ctx.ShouldBindJSON(&body); err != nil {
			logf(ctx, "login bind failed: %v", err)
			ctx.JSON(400, gin.H{
				"message": "Invalid request body",
			})
			return
		}
		if err := c.validate.Struct(body); err != nil {
			logf(ctx, "login field validation failed: %v", err)
			ctx.JSON(400, gin.H{
				"message": "validation failed, username and password are required!",
			})
			return
		}

		tokenString, err := c.svc.Login(service.LoginRequest{
			Username:  body.Username,
			UserPhone: body.UserPhone,
			Password:  body.Password,
		})
		if err != nil {
			logf(ctx, "login logic failed: %v", err)
			ctx.JSON(401, gin.H{
				"message": "wrong username or password",
			})
			return
		}
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

func logf(ctx *gin.Context, format string, args ...any) {
	if v, ok := ctx.Get("logf"); ok {
		if f, ok := v.(func(string, ...any)); ok {
			f(format, args...)
			return
		}
	}
	log.Printf("[authentication]"+format, args...)
}

func getValidationErrorMessage(v validator.FieldError) string {
	switch v.Tag() {
	case "required":
		return fmt.Sprintf("%s field must not be empty", v.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", v.Field())
	default:
		return fmt.Sprintf("%s field must not be empty", v.Field())
	}
}

func (c *Controller) HealthChecker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "healthy",
			"message": "your application is healthy",
		})
	}
}


