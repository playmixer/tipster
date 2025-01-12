package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/playmixer/tipster/internal/adapters/models"
	"go.uber.org/zap"
)

var (
	nameAccessToken  = "token"
	nameRefreshToken = "refresh"

	lifeTimeAccessToken int64 = 60 * 60 // seconds, 60 = 1 min
)

type tipster interface {
	SendOTP(ctx context.Context, email string) error
	CheckOTP(ctx context.Context, email string, otp string) (*models.User, error)
	NewRefreshToken(ctx context.Context, userID uint) (string, error)
	UpdRefreshToken(ctx context.Context, userID uint, refreshToken string) (string, error)
	DelRefreshToken(ctx context.Context, userID uint, refreshToken string) error

	GetUserByID(ctx context.Context, userID uint) (*models.User, error)

	Recognize(ctx context.Context, userID uint, fID string, data []byte, language string) (string, error)
	Translate(ctx context.Context, userID uint, fID string, sourceLanguage, targetLanguage string, text string) (string, error)
	Speech(ctx context.Context, userID uint, fID string, text, lang string) ([]byte, error)
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
	r.MaxMultipartMemory = 1 << 20
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")
	r.Use(
		s.Logger(),
	)

	r.GET("/logout", s.handlerLogout)

	auth := r.Group("/api/v0/auth")
	{
		auth.POST("/otp", s.handlerSendOTP)
		auth.POST("/signIn", s.handlerSignIn)
	}

	api := r.Group("/api/v0")
	api.Use(s.Authenticate())
	{
		api.POST("/audio/recognize", s.handlerAudioRecognize)
		api.POST("/text/translate", s.handlerTextTranslate)
		api.POST("/text/speech", s.handlerTextSpeech)
		api.GET("/info", s.handlerInfo)
	}

	page := r.Group("/")
	page.Use(s.Authenticate())
	{
		page.GET("/", func(c *gin.Context) {
			var email string
			if _email := s.extractJWT(c, "user_email"); _email != nil {
				if _email, ok := _email.(string); ok {
					email = _email
				}
			}
			c.HTML(http.StatusOK, "index.html", struct {
				Email string
			}{
				Email: email,
			})
			c.Errors.Errors()
		})
	}

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
		return fmt.Errorf("failed shutdown servr: %w", err)
	}
	return nil
}

func (s *Server) failed(c *gin.Context, message string, err error) {
	s.log.Error(message, zap.Error(err))
	c.Writer.WriteHeader(http.StatusInternalServerError)
}
