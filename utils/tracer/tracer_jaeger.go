package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Options struct {
	EndpointURL string
	ServiceName string
}

var singletonJaegerTracer = make(map[string]*tracerProviderJaeger)

type tracerProviderJaeger struct {
	provider *tracesdk.TracerProvider
	Error    error
}

func newJaegerTracer(tracerIdentifier string) *tracerProviderJaeger {
	provider := tracerProviderJaeger{}
	singletonJaegerTracer[tracerIdentifier] = provider.init(tracerIdentifier)
	return singletonJaegerTracer[tracerIdentifier]
}

func (tracer *tracerProviderJaeger) init(tracerIdentifier string) *tracerProviderJaeger {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithUsername(""),
		jaeger.WithPassword(""),
	))
	if err != nil {
		tracer.Error = err
		return nil
	}
	tracer.provider = tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(tracerIdentifier),
			attribute.String("environment", "development"),
		)),
	)

	//Register Tracer Provider
	otel.SetTracerProvider(tracer.provider)

	return tracer
}

func (tracer *tracerProviderJaeger) Finish(ctx context.Context, identifier string) {
	tracer.provider.Shutdown(ctx)
}
