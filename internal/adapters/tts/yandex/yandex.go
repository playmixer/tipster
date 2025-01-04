package yandex

import (
	"github.com/playmixer/tipster/internal/adapters/apperrors"
	"github.com/playmixer/tipster/pkg/yandex/cloud/tts"
)

type Config struct {
	APIKey string `env:"YANDEX_API_KEY"`
}

type Yandex struct {
	client *tts.TTS
}

func New(cfg Config) *Yandex {
	return &Yandex{
		client: tts.New(cfg.APIKey),
	}
}

func (y *Yandex) Speech(text string, lang string) ([]byte, error) {
	if l := tts.ParseLang(lang); l != nil {
		return y.client.Speech(text, tts.SetLanguage(*l))
	}

	return nil, apperrors.ErrLanguageNotSupportVoice
}
