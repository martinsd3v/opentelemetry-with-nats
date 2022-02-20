package usecase

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/cache"
	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/mysql"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type useCase struct {
	trc provider.Tracer
}

func New(trc provider.Tracer) *useCase {
	return &useCase{trc}
}

func (m *useCase) AuthUser(ctx context.Context, email, password string) bool {
	// ctx, span := m.trc.NewTracer(ctx, "UseCase").Span(ctx, "usecases/AuthUser")
	// defer span.Finish()

	repository := mysql.New(m.trc)
	repository.FindByEmail(ctx)
	repository.Insert(ctx)

	redis := cache.New(m.trc)
	redis.Get(ctx)

	return email == "email@gmail.com" && password == "password"
}

func (m *useCase) HashPassword(ctx context.Context) {
	// ctx, span := m.trc.NewTracer(ctx, "UseCase").Span(ctx, "usecases/HashPassword")
	// defer span.Finish()

	repository := mysql.New(m.trc)
	repository.Update(ctx)
}
