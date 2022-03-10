package cache

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

type cache struct{}

func New() *cache {
	return &cache{}
}

func (m *cache) Get(ctx context.Context) {
	_, span := tracer.Span(ctx, "Memcache/Get")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}

func (m *cache) Set(ctx context.Context) {
	_, span := tracer.Span(ctx, "Memcache/Set")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 50)))
}
