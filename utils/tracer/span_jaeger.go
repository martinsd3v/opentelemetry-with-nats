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

type SpanStartOption struct {
	Key   string
	Value interface{}
}

type Tracer struct {
	identifier string
	tracing    *tracing
	context    context.Context
	tracer     trace.Tracer
}

func (t Tracer) NewTracer(ctx context.Context, identifier string) Tracer {
	t.context = ctx
	t.identifier = identifier
	t.tracing = t.tracing.init(t.identifier)
	return t
}

func (t Tracer) Finish() {
	t.tracing.Finish(t.context, t.identifier)
}

func (t *Tracer) Span(ctx context.Context, identifier string, opts ...SpanStartOption) (context.Context, Span) {
	tr := otel.Tracer(t.identifier)
	options := t.parseOptions(opts...)
	ctx, spn := tr.Start(ctx, identifier, options...)
	return ctx, Span{spn}
}

func (t Tracer) ExportSpanContext() SpanContext {
	return SpanContext{
		TraceID:    trace.SpanContextFromContext(t.context).TraceID().String(),
		SpanID:     trace.SpanContextFromContext(t.context).SpanID().String(),
		Remote:     trace.SpanContextFromContext(t.context).IsRemote(),
		TraceState: trace.SpanContextFromContext(t.context).TraceState().String(),
		TraceFlags: byte(trace.SpanContextFromContext(t.context).TraceFlags()),
	}
}

func (t Tracer) parseOptions(startOptions ...SpanStartOption) []trace.SpanStartOption {
	opts := make([]trace.SpanStartOption, len(startOptions))

	for i, opt := range startOptions {
		value, _ := json.Marshal(opt.Value)
		attr := attribute.String(opt.Key, string(value))
		opts[i] = trace.WithAttributes(attr)
	}

	return opts
}

type Span struct {
	span trace.Span
}

func (s Span) Finish() {
	s.span.End()
}
