package fake

import "context"

type Fake struct{}

func New() *Fake {
	return &Fake{}
}

func (f *Fake) Set(ctx context.Context, k string, v any, lifeTime int64) {
	return
}
func (f *Fake) Get(ctx context.Context, k string) any {
	return nil
}
