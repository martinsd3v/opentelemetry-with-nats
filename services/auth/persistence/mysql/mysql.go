package mysql

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type mysql struct {
	trc provider.Tracer
}

func New(trc provider.Tracer) *mysql {
	trc = trc.New("Mysql")
	return &mysql{trc}
}

func (m *mysql) FindByEmail(ctx context.Context) {
	_, span := m.trc.Span(ctx, "Mysql/FindByEmail")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}

func (m *mysql) Insert(ctx context.Context) {
	_, span := m.trc.Span(ctx, "Mysql/Insert")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}

func (m *mysql) Update(ctx context.Context) {
	_, span := m.trc.Span(ctx, "Mysql/Update")
	defer span.End()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}
