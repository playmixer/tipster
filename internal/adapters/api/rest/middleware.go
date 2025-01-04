package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		s.log.Info(
			"Request",
			zap.String("uri", c.Request.RequestURI),
			zap.Duration("duration", time.Since(start)),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
