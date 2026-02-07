package controllers

import (
	"authentication/internal/service"
	"fmt"
	"log"
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sign up failed",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"user":    user.Username,
	})

}

func LoginGin(c *gin.Context, authService *service.AuthService) {
	fmt.Println("Login endpoint hit")
	var body service.LoginInfo
	if err := c.ShouldBindJSON(&body); err != nil || body.Username == "" || body.Password == "" {
		c.JSON(400, gin.H{
			"error": "Invalid request, username and password are required",
		})
		return
	}
	user, err := authService.Login(body)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"error": "wrong username or password",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "login successful",
		"user":    user.Username,
	})
}
