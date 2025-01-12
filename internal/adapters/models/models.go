package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"index:idx_email"`
	PasswordHash string
}

type Session struct {
	gorm.Model
	Token     string
	User      User
	UserID    uint
	ExpiresAt int64
}

type UserOTP struct {
	gorm.Model
	User      User
	UserID    uint
	Code      string
	ExpiresAt int64
}

type Recognize struct {
	gorm.Model
	User     User
	UserID   uint
	Range    float32
	FrontID  string
	Text     string
	Language string
}

type Translate struct {
	gorm.Model
	User         User
	UserID       uint
	FromText     string
	FromLanguage string
	ToText       string
	ToLanguage   string
	FrontID      string
}

type Speech struct {
	gorm.Model
	User     User
	UserID   uint
	Text     string
	Language string
	FrontID  string
}
