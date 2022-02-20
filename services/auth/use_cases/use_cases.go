package usecase

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/cache"
	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/persistence/mysql"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

type useCase struct{}

func New() *useCase {
	return &useCase{}
}

func (m *useCase) AuthUser(ctx context.Context, email, password string) bool {
	ctx, span := provider.Span(ctx, "usecases/AuthUser")
	defer span.End()

	repository := mysql.New()
	repository.FindByEmail(ctx)
	repository.Insert(ctx)

	redis := cache.New()
	redis.Get(ctx)

	return email == "email@gmail.com" && password == "password"
}

func (m *useCase) HashPassword(ctx context.Context) {
	ctx, span := provider.Span(ctx, "usecases/HashPassword")
	defer span.End()

	repository := mysql.New()
	repository.Update(ctx)
}
