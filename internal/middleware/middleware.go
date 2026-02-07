package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("requestID", requestID)
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		endpoint := c.Request.Method + " " + c.FullPath()
		log.Printf("[%s], RequestID: %s, latency: %s, status: %d, endpoint: %s", time.Now().Format(time.RFC3339), requestID, latency, status, endpoint)
	}
}

func JwtMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokeStr := c.GetHeader("Authorization")[7:]
		token, err := jwt.Parse(tokeStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
