package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func ContextFromSpanContext(ctx context.Context, spanContext *SpanContext) context.Context {
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
