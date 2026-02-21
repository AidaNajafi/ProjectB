package controllers

import (
	"authentication/internal/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type logLevel string

var (
	Info  logLevel = "INFO"
	Warn  logLevel = "WARN"
	Error logLevel = "ERROR"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx.Set("request_id", requestID)
		ctx.Writer.Header().Set("X-Request-ID", requestID)

		ctx.Set("logf", func(format string, args ...any) {
			log.Printf("authentication RequestID=%s"+format, append([]any{requestID}, args...)...)
		})

		start := time.Now()
		ctx.Next()
		latency := time.Since(start)
		status := ctx.Writer.Status()
		endpoint := ctx.FullPath()
		var level logLevel
		switch {
		case status >= 400:
			level = Error
		default:
			level = Info
		}

		log.Printf("[%s] time=%s, RequestID: %s, latency: %s, status: %d, endpoint: %s", level, time.Now().Format(time.RFC3339), requestID, latency, status, endpoint)
	}
}

func JwtMiddleware(jwtSvc *service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if len(auth) < 8 || auth[:7] != "Bearer " {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
			})
			return
		}
		tokenStr := auth[7:]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSvc.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			log.Println(err)
			ctx.JSON(401, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
