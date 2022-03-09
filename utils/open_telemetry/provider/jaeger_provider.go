package provider

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
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
	provider *tracesdk.TracerProvider
	Err      error
}

func (trc Tracer) Finish() {
	fmt.Println("Finalizou aqui")
	trc.provider.Shutdown(context.Background())
}

func (trc Tracer) New(serviceName string) Tracer {
	trc.opts.service = serviceName
	// trc.provider, trc.Err = initJeagerProvider(trc.opts)
	trc.provider, trc.Err = initSignozProvider(trc.opts)
	if trc.Err != nil {
		return trc
	}
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

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

func initSignozProvider(opts Options) (*tracesdk.TracerProvider, error) {
	headers := map[string]string{
		"signoz-access-token": "",
	}

	secureOption := otlptracegrpc.WithInsecure()
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(opts.EndpointURL),
			otlptracegrpc.WithHeaders(headers),
		),
	)

	if err != nil {
		return nil, err
	}

	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSpanProcessor(tracesdk.NewBatchSpanProcessor(exporter)),
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("service.name", opts.service),
			attribute.String("library.language", "go"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}
