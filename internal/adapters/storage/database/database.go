package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"github.com/playmixer/tipster/internal/adapters/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	lifeTimeOTP          int64 = 60 * 5
	lifeTimeRefreshToken int64 = 60 * 60 * 24 * 30
)

type Config struct {
	Name string `env:"DB_NAME"`
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
}

type Database struct {
	db  *gorm.DB
	log *zap.Logger
}

func New(cfg Config, log *zap.Logger) (*Database, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Europe/Moscow", cfg.User, cfg.Pass, cfg.Name, cfg.Host, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed open database: %w", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.UserOTP{},
		&models.Session{},
		&models.Recognize{},
		&models.Translate{},
		&models.Speech{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed auto migrate: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func (s *Database) NewOTP(ctx context.Context, email string, otp string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Email: email,
		}
		err := tx.Where("email = ?", email).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed find user by email `%s` :%w", email, err)
		}

		if user.ID == 0 {
			err = tx.Create(&user).Error
			if err != nil {
				return fmt.Errorf("failed update user: %w", err)
			}
		}

		err = tx.Where("user_id = ?", user.ID).Delete(&models.UserOTP{}).Error
		if err != nil {
			return fmt.Errorf("failed delete older otp: %w", err)
		}

		userOTP := models.UserOTP{
			UserID:    user.ID,
			Code:      otp,
			ExpiresAt: time.Now().UTC().Unix() + lifeTimeOTP,
		}
		err = tx.Save(&userOTP).Error
		if err != nil {
			return fmt.Errorf("failed save otp to store: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed save otp: %w", err)
	}
	return nil
}
func (s *Database) GetOTP(ctx context.Context, email string) (*models.UserOTP, error) {
	otp := &models.UserOTP{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user := &models.User{}
		err := tx.Where("email = ?", email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed get user by email: %w", err)
		}

		err = tx.Where("user_id = ?", user.ID).Preload("User").First(&otp).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed get otp by user id: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed find otp: %w", err)
	}

	return otp, nil
}

func (s *Database) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	user := &models.User{}
	err := s.db.Where("id = ?", userID).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrRecordNotFound
		}
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Database) NewRefreshToken(ctx context.Context, userID uint, token string, lifeTime int64) error {
	if lifeTime == 0 {
		lifeTime = lifeTimeRefreshToken
	}
	sess := &models.Session{UserID: userID, Token: token, ExpiresAt: time.Now().Unix() + lifeTime}
	err := s.db.Create(sess).Error
	if err != nil {
		return fmt.Errorf("failed create session: %w", err)
	}
	return nil
}

func (s *Database) UpdRefreshToken(ctx context.Context, userID uint, oldToken, newToken string, lifeTime int64) error {
	if lifeTime == 0 {
		lifeTime = lifeTimeRefreshToken
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		sess := &models.Session{}
		if err := tx.Where("token = ? and user_id = ?", oldToken, userID).First(sess).Error; err != nil {
			return fmt.Errorf("failed find session: %w", err)
		}

		if sess.ExpiresAt < time.Now().Unix() {
			return fmt.Errorf("refresh token expirated")
		}

		sess.Token = newToken
		sess.ExpiresAt = time.Now().Unix() + lifeTime
		if err := tx.Save(sess).Error; err != nil {
			return fmt.Errorf("failed update session: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed update refresh token: %w", err)
	}

	return nil
}

func (s *Database) DelRefreshToken(ctx context.Context, userID uint, refreshToken string) error {
	sess := &models.Session{}
	err := s.db.Where("user_id = ? and token = ?", userID, refreshToken).Delete(sess).Error
	if err != nil {
		return fmt.Errorf("failed delete refresh token: %w", err)
	}

	return nil
}

func (s *Database) NewRecognize(ctx context.Context, userID uint, fID string, _range float32, language string, text string) (*models.Recognize, error) {
	rec := &models.Recognize{
		UserID:   userID,
		Range:    _range,
		FrontID:  fID,
		Text:     text,
		Language: language,
	}

	err := s.db.Create(&rec).Error
	if err != nil {
		return nil, fmt.Errorf("failed create recognize meta: %w", err)
	}

	return rec, nil
}

func (s *Database) NewTranslate(ctx context.Context, userID uint, fID string, fromText, fromLanguage, toText, toLanguage string) (*models.Translate, error) {
	transl := &models.Translate{
		UserID:       userID,
		FromText:     fromText,
		FromLanguage: fromLanguage,
		ToText:       toText,
		ToLanguage:   toLanguage,
		FrontID:      fID,
	}
	err := s.db.Create(&transl).Error
	if err != nil {
		return nil, fmt.Errorf("failed create translate meta: %w", err)
	}

	return transl, nil
}

func (s *Database) NewSpeech(ctx context.Context, userID uint, fID string, text, language string) (*models.Speech, error) {
	speech := &models.Speech{
		UserID:   userID,
		Text:     text,
		Language: language,
		FrontID:  fID,
	}
	err := s.db.Create(&speech).Error
	if err != nil {
		return nil, fmt.Errorf("failed create speech meta: %w", err)
	}

	return speech, nil
}
