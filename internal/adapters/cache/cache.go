package cache

import (
	"context"

	"github.com/playmixer/tipster/internal/adapters/cache/fake"
	"github.com/playmixer/tipster/internal/adapters/cache/memory"
	"github.com/playmixer/tipster/internal/adapters/cache/redisstore"
	"go.uber.org/zap"
)

type Config struct {
	Redis redisstore.Config
}

type Cache interface {
	Set(ctx context.Context, k string, v any, lifeTime int64)
	Get(ctx context.Context, k string) any
}

func New(ctx context.Context, name string, cfg Config, log *zap.Logger) Cache {
	if name == "redis" {
		m := redisstore.New(cfg.Redis, log)
		log.Info("Initialize redis cache storage")
		return m
	}
	if name == "fake" {
		log.Info("Initialize fake cache storage")
		return fake.New()
	}
	log.Info("Initialize memory cache storage")
	m := memory.New(ctx)

	return m
}
