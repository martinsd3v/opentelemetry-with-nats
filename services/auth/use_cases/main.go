package usecase

import (
	"context"
	"opentelemetry/services/auth/persistence/cache"
	"opentelemetry/services/auth/persistence/mysql"
	"opentelemetry/utils/tracer"
)

type useCase struct{}

func New() *useCase {
	return &useCase{}
}

func (m *useCase) AuthUser(ctx context.Context, email, password string) bool {
	ctx, span := tracer.New(ctx).WithNewTrace("ServiceAuth", "usecases/AuthUser")
	defer span.Finish()

	repository := mysql.New()
	repository.FindByEmail(ctx)
	repository.Insert(ctx)

	redis := cache.New()
	redis.Get(ctx)

	return email == "email@gmail.com" && password == "password"
}

func (m *useCase) HashPassword(ctx context.Context) {
	ctx, span := tracer.New(ctx).WithNewTrace("ServiceAuth", "usecases/HashPassword")
	defer span.Finish()

	repository := mysql.New()
	repository.Update(ctx)
}
