package middleware

import (
	"time"

	"github.com/scrawld/library/zaplog"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		// Process request
		ctx.Next()
		// Stop timer
		end := time.Now()

		if v, exists := ctx.Get("logger"); exists {
			l, ok := v.(*zaplog.TracingLogger)
			if ok && l != nil {
				l.Infof("code: %d, take: %s", ctx.Writer.Status(), end.Sub(start))
			}
		}
	}
}
