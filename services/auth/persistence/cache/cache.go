package cache

import (
	"context"
	"opentelemetry/utils"
	"opentelemetry/utils/tracer"
	"time"
)

type cache struct {
}

func New() *cache {
	return &cache{}
}

func (m *cache) Get(ctx context.Context) {
	_, span := tracer.New(ctx).WithNewTrace("Cache", "Memcache/Get")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}

func (m *cache) Set(ctx context.Context) {
	_, span := tracer.New(ctx).WithNewTrace("Cache", "Memcache/Set")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}
