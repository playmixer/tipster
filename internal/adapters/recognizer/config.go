package recognizer

import (
	"github.com/playmixer/tipster/internal/adapters/recognizer/multi"
	"github.com/playmixer/tipster/internal/adapters/recognizer/vosk"
	"github.com/playmixer/tipster/internal/adapters/recognizer/yandex"
)

type Config struct {
	Yandex yandex.Config
	Vosk   vosk.Config
	Multi  multi.Config
}
