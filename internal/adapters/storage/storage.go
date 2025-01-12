package storage

import (
	"context"
	"fmt"

	"github.com/playmixer/tipster/internal/adapters/models"
	"github.com/playmixer/tipster/internal/adapters/storage/database"
	"go.uber.org/zap"
)

type Config struct {
	Database database.Config
}

type StorageI interface {
	NewOTP(ctx context.Context, email string, otp string) error
	GetOTP(ctx context.Context, email string) (*models.UserOTP, error)
	GetUser(ctx context.Context, userID uint) (*models.User, error)
	NewRefreshToken(ctx context.Context, userID uint, token string, lifeTime int64) error
	UpdRefreshToken(ctx context.Context, userID uint, oldToken, newToken string, lifeTime int64) error
	DelRefreshToken(ctx context.Context, userID uint, refreshToken string) error

	NewRecognize(ctx context.Context, userID uint, fID string, _range float32, language string, text string) (*models.Recognize, error)
	NewTranslate(ctx context.Context, userID uint, fID string, fromText, fromLanguage, toText, toLanguage string) (*models.Translate, error)
	NewSpeech(ctx context.Context, userID uint, fID string, text, language string) (*models.Speech, error)
}

func New(cfg Config, log *zap.Logger) (StorageI, error) {
	db, err := database.New(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("failed initialize database: %w", err)
	}

	return db, nil
}
