package logger

import (
	"log"

	"github.com/gin-gonic/gin"
)

func logf(ctx *gin.Context, format string, args ...any) {
	if v, ok := ctx.Get("logf"); ok {
		if f, ok := v.(func(string, ...any)); ok {
			f(format, args...)
			return
		}
	}
	log.Printf("[authentication]"+format, args...)
}
