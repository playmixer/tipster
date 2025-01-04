package memory

import (
	"context"
	"sync"
	"time"
)

var (
	defaultLifeTime int64 = 60
)

type item struct {
	data    interface{}
	expired int64
}

type Memory struct {
	items map[string]item
	mu    *sync.Mutex
}

func New(ctx context.Context) *Memory {
	m := &Memory{
		mu: &sync.Mutex{},
	}
	m.items = make(map[string]item)
	go m.gb(ctx)
	return m
}

func (m *Memory) Set(ctx context.Context, k string, v any, lifeTime int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if lifeTime == 0 {
		lifeTime = int64(defaultLifeTime)
	}

	m.items[k] = item{
		data:    v,
		expired: time.Now().Unix() + lifeTime,
	}
}

func (m *Memory) Get(ctx context.Context, k string) any {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.items[k]; ok {
		return v.data
	}

	return nil
}

func (m *Memory) gb(ctx context.Context) {
	tick := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			current := time.Now().Unix()
			for k, v := range m.items {
				if v.expired < current {
					m.mu.Lock()
					delete(m.items, k)
					m.mu.Unlock()
				}
			}
		}
	}
}
