package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/playmixer/tipster/internal/adapters/cache"
	"github.com/playmixer/tipster/internal/adapters/recognizer"
	yndxRecognize "github.com/playmixer/tipster/internal/adapters/recognizer/yandex"
	"github.com/playmixer/tipster/internal/adapters/translator"
	yndxTranslator "github.com/playmixer/tipster/internal/adapters/translator/yandex"
	"github.com/playmixer/tipster/internal/adapters/tts"
	"github.com/playmixer/tipster/internal/adapters/tts/yandex"
)

type Config struct {
	LogLVL         string `env:"LOG_LEVEL"`
	RecognizerName string `env:"RECOGNIZER_NAME"`
	Recognizer     recognizer.Config
	CacheName      string `env:"CACHE_NAME"`
	Cache          cache.Config
	TTSName        string `env:"TTS_NAME"`
	TTS            tts.Config
	Translator     translator.Config
	Address        string `env:"HTTP_ADDRESS"`
}

func Init() (*Config, error) {
	cfg := &Config{
		Recognizer: recognizer.Config{
			Yandex: yndxRecognize.Config{},
		},
		Translator: translator.Config{
			Yandex: yndxTranslator.Config{},
		},
		TTS: tts.Config{
			Yandex: yandex.Config{},
		},
		Cache: cache.Config{},
	}

	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return cfg, fmt.Errorf("failed load enviorements from file: %w", err)
	}

	if err := env.Parse(cfg); err != nil {
		return cfg, fmt.Errorf("failed parse env: %w", err)
	}

	return cfg, nil
}
