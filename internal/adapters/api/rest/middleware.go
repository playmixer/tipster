package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/playmixer/tipster/internal/adapters/apperrors"
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

func (s *Server) Authenticate() gin.HandlerFunc {
	toAuthorization := func(c *gin.Context) {
		c.HTML(http.StatusUnauthorized, "authorization.html", http.NoBody)
		c.Abort()
	}

	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			toAuthorization(c)
			return
		}

		if err = s.validJWT(tokenString); err != nil {
			if !errors.Is(err, apperrors.AccessTokenExpired) {
				toAuthorization(c)
				return
			}
			err = s.refreshJWT(c)
			if err != nil {
				s.log.Error("failed refresh token", zap.Error(err))
				toAuthorization(c)
				return
			}
		}

		c.Next()
	}
}
