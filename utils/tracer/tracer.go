package tracer

import (
	"context"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type SigNoz struct {
	Enabled      bool   `json:"enabled"`
	Token        string `json:"token"`
	CollectorURL string `json:"collectorURL"`
}

type Jaeger struct {
	Enabled      bool   `json:"enabled"`
	CollectorURL string `json:"collectorURL"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type Options struct {
	SigNoz      SigNoz `json:"sigNoz"`
	Jaeger      Jaeger `json:"jaeger"`
	Environment string `json:"environment"`
	service     string
}

type Tracer struct {
	opts     Options
	provider *tracesdk.TracerProvider
	Err      error
}

func (trc Tracer) Finish() {
	trc.provider.Shutdown(context.Background())
}

func (trc Tracer) New(serviceName string) Tracer {
	trc.opts.service = serviceName

	if trc.opts.SigNoz.Enabled {
		trc.provider, trc.Err = setupSigNozCollector(trc.opts)
	}

	if trc.opts.Jaeger.Enabled {
		trc.provider, trc.Err = setupJaegerCollector(trc.opts)
	}

	return trc
}

func Start(opts Options, serviceName string) Tracer {
	trc := Tracer{opts: opts}
	return trc.New(serviceName)
}
