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
	EndpointURL string `json:"endpointUrl"`
	AgentHost   string `json:"agentHost"`
	AgentPort   string `json:"agentPort"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	identifier  string
}

type Tracing struct {
	provider *tracesdk.TracerProvider
	options  Options
	Error    error
}

func SetupJeagerTracer(options Options) Tracing {
	trc := Tracing{options: options}
	trc.init(options.identifier)
	return trc
}

func (tracing Tracing) New(ctx context.Context) *span {
	return &span{
		context: ctx,
		tracing: tracing,
	}
}

func (tracing *Tracing) init(tracerIdentifier string) *Tracing {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithUsername(""),
		jaeger.WithPassword(""),
	))
	if err != nil {
		tracing.Error = err
		return nil
	}
	tracing.provider = tracesdk.NewTracerProvider(
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
	otel.SetTracerProvider(tracing.provider)

	return tracing
}

func (tracing Tracing) finish(ctx context.Context, identifier string) {
	tracing.provider.Shutdown(ctx)
}
