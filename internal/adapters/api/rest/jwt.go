package rest

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"github.com/playmixer/tipster/internal/adapters/models"
)

func (s *Server) newJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"user_email": user.Email,
		"expired":    time.Now().Unix() + lifeTimeAccessToken,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed create token: %w", err)
	}

	return tokenString, nil
}

func (s *Server) validJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return fmt.Errorf("failed read token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if _expired, ok := claims["expired"]; ok {
			if expired, ok := _expired.(float64); ok {
				if int64(expired) < time.Now().Unix() {
					return apperrors.AccessTokenExpired
				} else {
					return nil
				}
			}
		}
	}

	return errors.New("token not valid")
}

func (s *Server) getJWTUser(c *gin.Context) (uint, error) {
	tokenString, err := c.Cookie(nameAccessToken)
	if err != nil {
		return 0, fmt.Errorf("not found")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed read token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if _userID, ok := claims["user_id"]; ok {
			if userID, ok := _userID.(float64); ok {
				return uint(userID), nil
			}
		}
	}

	return 0, errors.New("token not valid")
}

func (s *Server) extractJWT(c *gin.Context, field string) any {
	tokenString, err := c.Cookie(nameAccessToken)
	if err != nil {
		return nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if res, ok := claims[field]; ok {
			return res
		}
	}

	return nil
}

func (s *Server) refreshJWT(c *gin.Context) error {
	tokenString, err := c.Cookie(nameAccessToken)
	if err != nil {
		return fmt.Errorf("token not found")
	}

	userID, err := s.getJWTUser(c)
	if err != nil {
		return fmt.Errorf("failed get user from jwt: %w", err)
	}

	if userID == 0 {
		return fmt.Errorf("user id is empty")
	}

	refreshToken, err := c.Cookie(nameRefreshToken)
	if err != nil {
		return fmt.Errorf("refresh token not found")
	}

	user, err := s.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		return fmt.Errorf("failed find user: %w", err)
	}

	refreshToken, err = s.service.UpdRefreshToken(c.Request.Context(), userID, refreshToken)
	if err != nil {
		return fmt.Errorf("failed upd refresh token: %w", err)
	}

	tokenString, err = s.newJWT(user)
	c.SetCookie(nameAccessToken, tokenString, 0, "/", c.Request.Host, true, true)
	c.SetCookie(nameRefreshToken, refreshToken, 0, "/", c.Request.Host, true, true)

	return nil
}
