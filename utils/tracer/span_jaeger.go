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

type Span struct {
	provider        *tracerProviderJaeger
	traceIdentifier string
	span            trace.Span
	context         context.Context
}

func New(ctx context.Context) *Span {
	return &Span{
		context: ctx,
	}
}

func (s *Span) Trace(traceIdentifier string) *Span {
	s.traceIdentifier = traceIdentifier
	s.provider = newJaegerTracer(s.traceIdentifier)
	return s
}

func (s *Span) getOptions(startOptions ...SpanStartOption) []trace.SpanStartOption {
	opts := make([]trace.SpanStartOption, len(startOptions))

	for i, opt := range startOptions {
		value, _ := json.Marshal(opt.Value)
		attr := attribute.String(opt.Key, string(value))
		opts[i] = trace.WithAttributes(attr)
	}

	return opts
}

func (s *Span) Simple(identifier string, opts ...SpanStartOption) (context.Context, *Span) {
	tr := otel.Tracer(identifier)
	options := s.getOptions(opts...)
	s.context, s.span = tr.Start(s.context, identifier, options...)
	return s.context, s
}

func (s *Span) WithNewTrace(traceIdentifier, spanIdentifier string, opts ...SpanStartOption) (context.Context, *Span) {
	s.traceIdentifier = traceIdentifier
	s.provider = newJaegerTracer(traceIdentifier)

	tr := otel.Tracer(spanIdentifier)
	options := s.getOptions(opts...)
	s.context, s.span = tr.Start(s.context, spanIdentifier, options...)
	return s.context, s
}

func (s *Span) ExportSpanContext() SpanContext {
	return SpanContext{
		TraceID:    trace.SpanContextFromContext(s.context).TraceID().String(),
		SpanID:     trace.SpanContextFromContext(s.context).SpanID().String(),
		Remote:     trace.SpanContextFromContext(s.context).IsRemote(),
		TraceState: trace.SpanContextFromContext(s.context).TraceState().String(),
		TraceFlags: byte(trace.SpanContextFromContext(s.context).TraceFlags()),
	}
}

func (s *Span) Finish() {
	if s.span != nil {
		s.span.End()
	}
	if s.provider != nil {
		s.provider.Finish(s.context, s.traceIdentifier)
	}
}

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
