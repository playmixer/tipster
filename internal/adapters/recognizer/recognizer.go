package recognizer

import (
	"errors"
	"fmt"

	"github.com/playmixer/tipster/internal/adapters/recognizer/multi"
	"github.com/playmixer/tipster/internal/adapters/recognizer/vosk"
	"github.com/playmixer/tipster/internal/adapters/recognizer/yandex"
	"go.uber.org/zap"
)

type recognizerI interface {
	Recognize(data []byte, language string) (string, error)
}

func New(name string, cfg Config, log *zap.Logger) (recognizerI, error) {
	if name == "multi" {
		r, err := multi.New(cfg.Multi, log)
		if err != nil {
			return nil, fmt.Errorf("failed initizalize multi recognizer: %w", err)
		}

		return r, nil
	}

	if name == "vosk" {
		r := vosk.New(cfg.Vosk)
		return r, nil
	}

	if name == "yandex" {
		r, err := yandex.New(cfg.Yandex)
		if err != nil {
			return nil, fmt.Errorf("failed initizalize yandex recognizer: %w", err)
		}
		return r, nil
	}

	return nil, errors.New("failed found recognizaer")
}
