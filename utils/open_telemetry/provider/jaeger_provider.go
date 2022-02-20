package provider

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
	AgentConnect bool   `json:"agentConnect"`
	AgentHost    string `json:"agentHost"`
	AgentPort    string `json:"agentPort"`
	EndpointURL  string `json:"endpointUrl"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Environment  string `json:"environment"`
	service      string
}

type Tracer struct {
	opts     Options
	Provider *tracesdk.TracerProvider
	Err      error
}

func (trc Tracer) Finish() {
	trc.Provider.Shutdown(context.Background())
}

func (trc Tracer) New(serviceName string) Tracer {
	trc.opts.service = serviceName
	trc.Provider, trc.Err = initJeagerProvider(trc.opts)
	if trc.Err != nil {
		return trc
	}
	otel.SetTracerProvider(trc.Provider)
	return trc
}

func Start(opts Options, serviceName string) Tracer {
	trc := Tracer{opts: opts}
	return trc.New(serviceName)
}

func initJeagerProvider(opts Options) (*tracesdk.TracerProvider, error) {
	var exporter *jaeger.Exporter
	var err error

	if opts.AgentConnect {
		exporter, err = jaeger.New(jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(opts.AgentHost),
			jaeger.WithAgentPort(opts.AgentPort),
		))
	} else {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(opts.EndpointURL),
			jaeger.WithUsername(opts.Username),
			jaeger.WithPassword(opts.Username),
		))
	}

	if err != nil {
		return nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(opts.service),
			attribute.String("environment", opts.Environment),
		)),
	)

	return tracerProvider, nil
}
