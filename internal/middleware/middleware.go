package middleware

import (
	"authentication/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		ctx.Set("requestID", requestID)
		ctx.Next()
		latency := time.Since(start)
		status := ctx.Writer.Status()
		endpoint := ctx.Request.Method + " " + ctx.FullPath()
		log.Printf("[%s], RequestID: %s, latency: %s, status: %d, endpoint: %s", time.Now().Format(time.RFC3339), requestID, latency, status, endpoint)
	}
}

func JwtMiddleware(jwtSvc *service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokeStr := ctx.GetHeader("Authorization")[7:]
		if tokeStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
		}

		token, err := jwt.Parse(tokeStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSvc, nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
