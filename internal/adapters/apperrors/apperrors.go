package apperrors

import "errors"

var (
	ErrLanguageNotSupportVoice     = errors.New("voice is not supported for this language")
	ErrLanguageNotSupportTranslace = errors.New("language not supported for translation")
)
