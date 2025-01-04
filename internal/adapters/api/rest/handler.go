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

func (s *Server) handlerAudioRecognize(c *gin.Context) {
	f, err := c.FormFile("data")
	if err != nil {
		s.log.Error("failed getting audio", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	src, err := f.Open()
	if err != nil {
		s.log.Error("failed open audio", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(src)
	if err != nil {
		s.log.Error("failed read auido", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	language := c.PostForm("language")

	text, err := s.service.Recognize(c.Request.Context(), data, language)
	if err != nil {
		s.log.Error("failed recognize", zap.Error(err))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"language": language,
		"text":     text,
	})
}

func (s *Server) handlerTextTranslate(c *gin.Context) {
	req := textTranslateRequest{}

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

	text, err := s.service.Translate(c.Request.Context(), req.SourceLang, req.TargetLang, req.Text)
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

	speech, err := s.service.Speech(c.Request.Context(), req.Text, req.Lang)
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

	languages := map[string]string{
		"RU (русский)":       "ru-RU",
		"US (US)":            "en-US",
		"DE (DE)":            "de-DE",
		"ES (ES)":            "es-ES",
		"FI (FI)":            "fi-FI",
		"FR (FR)":            "fr-FR",
		"HE (HE)":            "he-HE",
		"IT (итальянский)":   "it-IT",
		"KZ (казахский)":     "kk-KZ",
		"NL (голландский)":   "nl-NL",
		"PL (польский)":      "pl-PL",
		"PT (португальский)": "pt-PT",
		"BR (бразильский португальский)": "pt-BR",
		"SE (SE)": "sv-SE",
		"TR (TR)": "tr-TR",
		"UZ (UZ)": "uz-UZ",
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}
