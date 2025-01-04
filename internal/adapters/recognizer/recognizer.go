package recognizer

import (
	"errors"
	"fmt"

	"github.com/playmixer/tipster/internal/adapters/recognizer/yandex"
)

type Recognizer interface {
	Recognize(data []byte, language string) (string, error)
}

func New(name string, cfg Config) (Recognizer, error) {
	if name == "yandex" {
		r, err := yandex.New(cfg.Yandex)
		if err != nil {
			return nil, fmt.Errorf("failed initizalize yandex recognizer: %w", err)
		}
		return r, nil
	}

	return nil, errors.New("failed found recognizaer")
}
