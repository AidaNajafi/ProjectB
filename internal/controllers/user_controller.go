package controllers

import (
	"authentication/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpGin(c *gin.Context, authService *service.AuthService) {

	var body service.SignUpInfo

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request",
		})
		return
	}

	user, err := authService.SignUp(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"user":    user.Username,
	})

}

func LoginGin(c *gin.Context, authService *service.AuthService) {
	var body service.LoginInfo
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}
	user, err := authService.Login(body)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "wrong credentials",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "login successful",
		"user":    user.Username,
	})
}
