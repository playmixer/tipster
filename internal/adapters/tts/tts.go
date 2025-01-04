package tts

import (
	"fmt"

	"github.com/playmixer/tipster/internal/adapters/tts/yandex"
)

type ttsI interface {
	Speech(text string, lang string) ([]byte, error)
}

type Config struct {
	Yandex yandex.Config
}

func New(name string, cfg Config) (ttsI, error) {
	if name == "yandex" {
		return yandex.New(cfg.Yandex), nil
	}

	return nil, fmt.Errorf("failed find tts")
}
