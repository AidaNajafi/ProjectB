package controllers

import (
	"authentication/internal/logger"
	"authentication/internal/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func GetRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx.Set(RequestIDKey, requestID)
		ctx.Writer.Header().Set("X-Request-ID", requestID)

		ctx.Next()

	}
}

func JwtMiddleware(jwtSvc *service.JwtService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" || token == header {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "missing bearer token",
			})
		}
		user, err := jwtSvc.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "invalid or expired token",
			})
		}
		ctx.Set("user_id", user.ID)
		ctx.Next()

	}
}

func GeneralLog(l logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		reqId, _ := ctx.Get(RequestIDKey)
		l.Info("http_request", map[string]any{
			"method":     ctx.Request.Method,
			"status":     ctx.Writer.Status(),
			"duration":   time.Since(start).String(),
			"request_id": reqId,
		})
	}
}
