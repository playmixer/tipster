package multi

import (
	"fmt"

	"github.com/playmixer/tipster/internal/adapters/recognizer/vosk"
	"github.com/playmixer/tipster/internal/adapters/recognizer/yandex"
	"go.uber.org/zap"
)

type Config struct {
	Yandex yandex.Config
	Vosk   vosk.Config
}

type multiRecognize struct {
	yandex *yandex.Yandex
	vosk   *vosk.Client
}

func New(cfg Config, log *zap.Logger) (*multiRecognize, error) {
	var err error
	r := &multiRecognize{}
	r.yandex, err = yandex.New(cfg.Yandex)
	if err != nil {
		return nil, fmt.Errorf("failed initizalize yandex recognizer: %w", err)
	}
	r.vosk = vosk.New(cfg.Vosk, vosk.SetSampleRate(cfg.Vosk.SampleRate), vosk.SetLogger(log))

	return r, nil
}

func (m *multiRecognize) Recognize(data []byte, language string) (string, error) {
	if language == "ru" || language == "ru-RU" {
		text, err := m.vosk.Recognize(data, "")
		if err != nil {
			return "", fmt.Errorf("failed recognize by vosk: %w", err)
		}
		return text, nil
	}
	text, err := m.yandex.Recognize(data, language)
	if err != nil {
		return "", fmt.Errorf("failed recognize by yandex: %w", err)
	}

	return text, nil
}
