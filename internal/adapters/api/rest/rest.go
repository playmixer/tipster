package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type tipster interface {
	Recognize(ctx context.Context, data []byte, language string) (string, error)
	Translate(ctx context.Context, sourceLanguage, targetLanguage string, text string) (string, error)
	Speech(ctx context.Context, lang, text string) ([]byte, error)
}

type Server struct {
	srv     *http.Server
	log     *zap.Logger
	service tipster
	secret  []byte
}

type Option func(*Server)

func Logger(log *zap.Logger) Option {
	return func(s *Server) {
		s.log = log
	}
}

func SetAddress(address string) Option {
	return func(s *Server) {
		s.srv.Addr = address
	}
}

func SetSecretKey(key []byte) Option {
	return func(s *Server) {
		s.secret = key
	}
}

func New(service tipster, options ...Option) (*Server, error) {
	s := &Server{
		srv:     &http.Server{},
		log:     zap.NewNop(),
		service: service,
	}

	for _, opt := range options {
		opt(s)
	}

	r := gin.New()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")
	r.Use(
		s.Logger(),
	)
	api := r.Group("/api/v0")
	{
		api.POST("/audio/recognize", s.handlerAudioRecognize)
		api.POST("/text/translate", s.handlerTextTranslate)
		api.POST("/text/speech", s.handlerTextSpeech)
		api.GET("/info", s.handlerInfo)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", http.NoBody)
	})

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.srv.Handler = r.Handler()

	return s, nil
}

func (s *Server) Run() error {
	if err := s.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server stopped with error: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed shutfdown servr: %w", err)
	}
	return nil
}
