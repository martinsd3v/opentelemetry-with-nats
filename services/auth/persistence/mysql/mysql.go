package mysql

import (
	"context"
	"time"

	"github.com/martinsd3v/opentelemetry-with-nats/utils"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

type mysql struct{}

func New() *mysql {
	return &mysql{}
}

func (m *mysql) FindByEmail(ctx context.Context) {
	_, span := tracer.New(ctx).WithNewTrace("Mysql", "Mysql/FindByEmail")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}

func (m *mysql) Insert(ctx context.Context) {
	_, span := tracer.New(ctx).WithNewTrace("Mysql", "Mysql/Insert")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}

func (m *mysql) Update(ctx context.Context) {
	_, span := tracer.New(ctx).WithNewTrace("Mysql", "Mysql/Update")
	defer span.Finish()

	time.Sleep(time.Millisecond * time.Duration(utils.RandNumber(10, 300)))
}
