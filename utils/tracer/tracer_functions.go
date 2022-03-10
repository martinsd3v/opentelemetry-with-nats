package tracer

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

func InjectContext(ctx context.Context, spanContext *SpanContext) context.Context {
	if spanContext == nil {
		return ctx
	}

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

func Span(ctx context.Context, identifier string, opts ...SpanStartOption) (context.Context, trace.Span) {
	options := parseOptions(opts...)
	return otel.Tracer(identifier).Start(ctx, identifier, options...)
}

type SpanStartOption struct {
	Key   string
	Value interface{}
}

func parseOptions(startOptions ...SpanStartOption) []trace.SpanStartOption {
	opts := make([]trace.SpanStartOption, len(startOptions))

	for i, opt := range startOptions {
		value, _ := json.Marshal(opt.Value)
		attr := attribute.String(opt.Key, string(value))
		opts[i] = trace.WithAttributes(attr)
	}

	return opts
}
