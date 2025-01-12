package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"go.uber.org/zap"
)

func (s *Server) handlerSendOTP(c *gin.Context) {
	bBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	req := sendOTPRequest{}
	err = json.Unmarshal(bBody, &req)
	if err != nil {
		s.failed(c, "failed unmarshal request body", err)
		return
	}

	err = s.service.SendOTP(c.Request.Context(), req.Email)
	if err != nil {
		s.failed(c, "failed send otp to email", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func (s *Server) handlerSignIn(c *gin.Context) {
	bBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	req := signInRequest{}
	err = json.Unmarshal(bBody, &req)
	if err != nil {
		s.failed(c, "failed unmarshal request body", err)
		return
	}

	user, err := s.service.CheckOTP(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrPassNotEqual) || errors.Is(err, apperrors.ErrUserNotFoundWithThisEmail) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Your login or password is incorrect. Please try again",
			})
			return
		}
		s.failed(c, "failed check otp", err)
		return
	}

	token, err := s.newJWT(user)
	if err != nil {
		s.failed(c, "failed create token", err)
		return
	}

	refreshToken, err := s.service.NewRefreshToken(c.Request.Context(), user.ID)
	if err != nil {
		s.failed(c, "failed create refresh token", err)
		return
	}

	c.SetCookie(nameAccessToken, token, 0, "/", c.Request.Host, true, true)
	c.SetCookie(nameRefreshToken, refreshToken, 0, "/", c.Request.Host, true, true)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func (s *Server) handlerLogout(c *gin.Context) {
	if refreshToken, err := c.Cookie(nameRefreshToken); err == nil {
		if userID, err := s.getJWTUser(c); err == nil {
			err = s.service.DelRefreshToken(c.Request.Context(), userID, refreshToken)
			if err != nil {
				s.log.Error("failed delete refresh token", zap.Error(err))
			}
		}
	}
	c.SetCookie(nameAccessToken, "", 0, "/", c.Request.Host, true, true)
	c.SetCookie(nameRefreshToken, "", 0, "/", c.Request.Host, true, true)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) handlerAudioRecognize(c *gin.Context) {
	userID, _ := s.getJWTUser(c)

	f, err := c.FormFile("data")
	if err != nil {
		s.log.Error("failed getting audio", zap.Error(err))
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	src, err := f.Open()
	if err != nil {
		s.failed(c, "failed open audio", err)
		return
	}
	data, err := io.ReadAll(src)
	if err != nil {
		s.failed(c, "failed read auido", err)
		return
	}

	language := c.PostForm("language")

	text, err := s.service.Recognize(c.Request.Context(), userID, c.PostForm("frontendID"), data, language)
	if err != nil {
		if errors.Is(err, apperrors.ErrShortRecord) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": apperrors.ErrShortRecord.Error(),
			})
			return
		}
		s.failed(c, "failed recognize", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"language": language,
		"text":     text,
	})
}

func (s *Server) handlerTextTranslate(c *gin.Context) {
	userID, _ := s.getJWTUser(c)
	req := textTranslateRequest{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.failed(c, "failed read body request", err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		s.failed(c, "failed unmarshal body request", err)
		return
	}

	text, err := s.service.Translate(c.Request.Context(), userID, req.FrontendID, req.SourceLang, req.TargetLang, req.Text)
	if err != nil {
		s.log.Error("failed translate",
			zap.String("text", req.Text),
			zap.String("source", req.SourceLang),
			zap.String("target", req.TargetLang),
			zap.Error(err),
		)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"text":   text,
	})
}

func (s *Server) handlerTextSpeech(c *gin.Context) {
	userID, _ := s.getJWTUser(c)
	req := textSpeechRequest{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.log.Error("failed read body request", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		s.log.Error("failed unmarshal body request", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	speech, err := s.service.Speech(c.Request.Context(), userID, req.FrontendID, req.Text, req.Lang)
	if err != nil {
		if errors.Is(err, apperrors.ErrLanguageNotSupportVoice) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": apperrors.ErrLanguageNotSupportVoice.Error(),
			})
			return
		}
		s.log.Error("failed speech", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// f, err := os.Create("./audio.wav")
	// if err != nil {
	// 	s.log.Error("failed open audio", zap.Error(err))
	// }
	// if err == nil {
	// 	_, err = f.Write(speech)
	// 	if err != nil {
	// 		s.log.Error("failed write audio", zap.Error(err))
	// 	}
	// }
	// f.Close()

	c.Data(http.StatusOK, "audio/ogg", speech)
}

func (s *Server) handlerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"languages": s.service.GetLanguages(c.Request.Context()),
		"recognize": map[string]any{
			"languages":     s.service.GetRecognizeLanguages(c.Request.Context()),
			"maximumLength": 10,
		},
		"speech": map[string]any{
			"languages": s.service.GetSpeechLanguages(c.Request.Context()),
		},
	})
}
