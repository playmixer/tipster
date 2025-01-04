package yandex

type Config struct {
	APIKey   string `env:"YANDEX_API_KEY"`
	FolderID string `env:"YANDEX_FOLDER"`
}
