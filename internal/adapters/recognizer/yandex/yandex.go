package yandex

import (
	"fmt"

	"github.com/playmixer/tipster/pkg/yandex/cloud/recognize"
)

type Yandex struct {
	client *recognize.Recognizer
}

func New(cfg Config) (*Yandex, error) {
	client := recognize.New(cfg.APIKey)
	return &Yandex{client: client}, nil
}

func (y *Yandex) Recognize(data []byte, language string) (string, error) {
	res, err := y.client.Recognize(data, recognize.SetLanguage(recognize.Language(language)), recognize.SetFormat(recognize.LPCM))
	if err != nil {
		return "", fmt.Errorf("failed recognize speech: %w", err)
	}

	return res, nil
}
