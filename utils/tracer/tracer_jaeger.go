package tracer

import (
	"context"
	"fmt"

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
}

type tracing struct {
	provider *tracesdk.TracerProvider
	options  Options
}

func New(options Options) Tracer {
	trc := tracing{options: options}
	return Tracer{tracing: &trc}
}

func (tr tracing) init(identifier string) *tracing {
	fmt.Println("Iniciou :" + identifier)

	fmt.Println("Tenta conectar UDP")
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(tr.options.AgentHost),
		jaeger.WithAgentPort(tr.options.AgentPort),
	))

	if err != nil {
		fmt.Println("Tenta conectar TCP")
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(tr.options.EndpointURL),
			jaeger.WithUsername(tr.options.Username),
			jaeger.WithPassword(tr.options.Username),
		))
		if err != nil {
			fmt.Println("Error Jaeger :", err)
			return nil
		}
	}

	tr.provider = tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(identifier),
			attribute.String("environment", "development"),
		)),
	)

	//Register Tracer Provider
	otel.SetTracerProvider(tr.provider)

	return &tr
}

func (tr tracing) Finish(ctx context.Context, identifier string) {
	fmt.Println("Finalizou :" + identifier)
	tr.provider.Shutdown(ctx)
}
