package cache

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

type cache struct {
	tracing tracer.Tracing
}

func New(tracing tracer.Tracing) *cache {
	return &cache{tracing}
}

func (m *cache) Get(ctx context.Context) {
	_, span := m.tracing.New(ctx).WithNewTrace("Cache", "Memcache/Get")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}

func (m *cache) Set(ctx context.Context) {
	_, span := m.tracing.New(ctx).WithNewTrace("Cache", "Memcache/Set")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}
