package usecase

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/cache"
	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/mysql"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

type useCase struct {
	tracing tracer.Tracing
}

func New(tracing tracer.Tracing) *useCase {
	return &useCase{tracing}
}

func (m *useCase) AuthUser(ctx context.Context, email, password string) bool {
	ctx, span := m.tracing.New(ctx).WithNewTrace("ServiceAuth", "usecases/AuthUser")
	defer span.Finish()

	repository := mysql.New(m.tracing)
	repository.FindByEmail(ctx)
	repository.Insert(ctx)

	redis := cache.New(m.tracing)
	redis.Get(ctx)

	return email == "email@gmail.com" && password == "password"
}

func (m *useCase) HashPassword(ctx context.Context) {
	ctx, span := m.tracing.New(ctx).WithNewTrace("ServiceAuth", "usecases/HashPassword")
	defer span.Finish()

	repository := mysql.New(m.tracing)
	repository.Update(ctx)
}
