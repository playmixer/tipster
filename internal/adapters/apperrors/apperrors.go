package apperrors

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")

	ErrLanguageNotSupportVoice     = errors.New("voice is not supported for this language")
	ErrLanguageNotSupportTranslace = errors.New("language not supported for translation")

	ErrShortRecord = errors.New("record is short")
	ErrLongRecord  = errors.New("record is long")

	//auth
	ErrPassNotEqual              = errors.New("passwords not equal")
	ErrUserNotFoundWithThisEmail = errors.New("user not found with this email")
	AccessTokenExpired           = errors.New("access token expired")
)
