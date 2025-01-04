package yandex

import (
	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"github.com/playmixer/tipster/pkg/yandex/cloud/translator"
)

type Config struct {
	APIKey string `env:"YANDEX_API_KEY"`
}

type Translator struct {
	client *translator.Translator
}

func New(cfg Config) *Translator {
	t := &Translator{
		client: translator.New(cfg.APIKey),
	}

	return t
}

func (t *Translator) Translate(sourceLang, targetLang string, text string) (string, error) {
	sourceLang = translator.ParseLang(sourceLang)
	targetLang = translator.ParseLang(targetLang)
	if sourceLang == "" || targetLang == "" {
		return "", apperrors.ErrLanguageNotSupportTranslace
	}

	return t.client.Translate(sourceLang, targetLang, text)
}
