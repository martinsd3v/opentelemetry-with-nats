package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type SpanContext struct {
	TraceID    string `json:"traceId"`
	SpanID     string `json:"spanId"`
	TraceState string `json:"traceState"`
	Remote     bool   `json:"remote"`
	TraceFlags byte   `json:"traceFlags"`
}

func ExportSpanContext(ctx context.Context) SpanContext {
	return SpanContext{
		TraceID:    trace.SpanContextFromContext(ctx).TraceID().String(),
		SpanID:     trace.SpanContextFromContext(ctx).SpanID().String(),
		Remote:     trace.SpanContextFromContext(ctx).IsRemote(),
		TraceState: trace.SpanContextFromContext(ctx).TraceState().String(),
		TraceFlags: byte(trace.SpanContextFromContext(ctx).TraceFlags()),
	}
}

func ImportContext(ctx context.Context, spanContext *SpanContext) context.Context {
	if spanContext != nil {
		traceID, _ := trace.TraceIDFromHex(spanContext.TraceID)
		spanID, _ := trace.SpanIDFromHex(spanContext.SpanID)
		traceState, _ := trace.ParseTraceState(spanContext.TraceState)
		newSpanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			Remote:     true,
			TraceState: traceState,
			TraceFlags: trace.TraceFlags(spanContext.TraceFlags),
		})

		return trace.ContextWithSpanContext(ctx, newSpanContext)
	}

	return ctx
}
