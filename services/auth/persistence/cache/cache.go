package cache

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type cache struct {
	trc provider.Tracer
}

func New(trc provider.Tracer) *cache {
	return &cache{trc}
}

func (m *cache) Get(ctx context.Context) {
	// _, span := m.trc.NewTracer(ctx, "Cache").Span(ctx, "Memcache/Get")
	// defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}

func (m *cache) Set(ctx context.Context) {
	// _, span := m.trc.NewTracer(ctx, "Cache").Span(ctx, "Memcache/Set")
	// defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}
