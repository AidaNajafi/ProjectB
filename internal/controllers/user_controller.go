package controllers

import (
	"authentication/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpGin(c *gin.Context) {

	var body service.SignUpInfo

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, err := service.SignUp(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "SignUP failed!",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "SignUP successful",
		"user":    user.Username,
	})

}
