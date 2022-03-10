package tracer

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func setupJaegerCollector(opts Options) (*tracesdk.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(opts.Jaeger.CollectorURL),
		jaeger.WithUsername(opts.Jaeger.Username),
		jaeger.WithPassword(opts.Jaeger.Password),
	))

	if err != nil {
		return nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("service.name", opts.service),
			attribute.String("library.language", "go"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}
