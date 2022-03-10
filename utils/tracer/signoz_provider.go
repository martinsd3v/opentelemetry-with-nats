package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func setupSigNozCollector(opts Options) (*tracesdk.TracerProvider, error) {
	headers := map[string]string{
		"signoz-access-token": opts.SigNoz.Token,
	}

	secureOption := otlptracegrpc.WithInsecure()
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(opts.SigNoz.CollectorURL),
			otlptracegrpc.WithHeaders(headers),
		),
	)

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
