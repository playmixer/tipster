package redisstore

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Store struct {
	store *redis.Client
	log   *zap.Logger
}

func New(cfg Config, log *zap.Logger) *Store {
	r := &Store{
		log: log,
	}

	r.store = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	return r
}

func (s *Store) Set(ctx context.Context, k string, v any, lifeTime int64) {
	err := s.store.Set(ctx, k, v, time.Duration(lifeTime)*time.Second).Err()
	if err != nil {
		s.log.Error("failed set to redis", zap.Error(err))
	}
}

func (s *Store) Get(ctx context.Context, k string) any {
	v, err := s.store.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		s.log.Error("failed get from redis", zap.Error(err))
		return nil
	}

	return v
}
