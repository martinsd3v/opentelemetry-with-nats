package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Span(ctx context.Context, identifier string) (context.Context, trace.Span) {
	return otel.Tracer(identifier).Start(ctx, identifier)
}
