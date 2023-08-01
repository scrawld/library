package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check for incoming header, use it if exists
		requestId := ctx.GetHeader("X-Request-Id")

		// Create request id with UUID4
		if requestId == "" {
			uid, _ := uuid.NewRandom()
			requestId = uid.String()
		}

		// Expose it for use in the application
		ctx.Request.Header.Set("X-Request-Id", requestId)

		// Set X-Request-Id header
		ctx.Header("X-Request-Id", requestId)
		ctx.Next()
	}
}
