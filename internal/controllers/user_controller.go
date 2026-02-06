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
		"message": "SignUP successful",
		"user":    user.Username,
	})

}
