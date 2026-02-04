package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(){
	router := gin.Default()
	router.POST("login", Login())
	router.POST("/signup", SignUp())
	
	protec
}



func Login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "login successful",
	})
}


func SignUp(c *gin.Context){
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "signup successful"
	})
}
