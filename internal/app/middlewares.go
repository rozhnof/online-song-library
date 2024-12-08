package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		if c.Writer.Status() >= http.StatusInternalServerError {
			logger.InfoContext(
				c.Request.Context(),
				"internal server error",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Int("status", c.Writer.Status()),
				slog.String("address", c.Request.RemoteAddr),
				slog.String("duration", duration.String()),
			)
		} else if duration > time.Second {
			logger.InfoContext(
				c.Request.Context(),
				"long time response",
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.Int("status", c.Writer.Status()),
				slog.String("address", c.Request.RemoteAddr),
				slog.String("duration", duration.String()),
			)
		}
	}
}
