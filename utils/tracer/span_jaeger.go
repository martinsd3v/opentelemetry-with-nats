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

type span struct {
	traceIdentifier string
	traceSpan       trace.Span
	context         context.Context
	tracing         Tracing
}

func (s *span) Trace(traceIdentifier string) *span {
	s.traceIdentifier = traceIdentifier
	s.tracing.init(s.traceIdentifier)
	return s
}

func (s *span) Simple(identifier string, opts ...SpanStartOption) (context.Context, spanCloser) {
	tr := otel.Tracer(identifier)
	options := s.parseOptions(opts...)
	s.context, s.traceSpan = tr.Start(s.context, identifier, options...)
	return s.context, spanCloser{span: *s}
}

func (s *span) WithNewTrace(traceIdentifier, spanIdentifier string, opts ...SpanStartOption) (context.Context, spanCloser) {
	s.traceIdentifier = traceIdentifier
	s.tracing.init(s.traceIdentifier)

	tr := otel.Tracer(spanIdentifier)
	options := s.parseOptions(opts...)
	s.context, s.traceSpan = tr.Start(s.context, spanIdentifier, options...)
	return s.context, spanCloser{span: *s}
}

func (s *span) ExportSpanContext() SpanContext {
	return SpanContext{
		TraceID:    trace.SpanContextFromContext(s.context).TraceID().String(),
		SpanID:     trace.SpanContextFromContext(s.context).SpanID().String(),
		Remote:     trace.SpanContextFromContext(s.context).IsRemote(),
		TraceState: trace.SpanContextFromContext(s.context).TraceState().String(),
		TraceFlags: byte(trace.SpanContextFromContext(s.context).TraceFlags()),
	}
}

func (s *span) parseOptions(startOptions ...SpanStartOption) []trace.SpanStartOption {
	opts := make([]trace.SpanStartOption, len(startOptions))

	for i, opt := range startOptions {
		value, _ := json.Marshal(opt.Value)
		attr := attribute.String(opt.Key, string(value))
		opts[i] = trace.WithAttributes(attr)
	}

	return opts
}

type spanCloser struct{ span }

func (s spanCloser) Finish() {
	if s.traceSpan != nil {
		s.traceSpan.End()
	}
	s.tracing.finish(s.context, s.traceIdentifier)
}
