package cache

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type cache struct {
	tracer provider.Tracer
}

func New(tracer provider.Tracer) *cache {
	tracer = tracer.New("Cache")
	return &cache{tracer}
}

func (m *cache) Get(ctx context.Context) {
	_, span := m.tracer.Span(ctx, "Memcache/Get")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}

func (m *cache) Set(ctx context.Context) {
	_, span := m.tracer.Span(ctx, "Memcache/Set")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}
