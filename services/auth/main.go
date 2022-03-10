package main

import (
	"fmt"
	"runtime"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

func main() {
	//Tracer
	trc := tracer.Start(tracer.Options{
		Environment: "dev",
		Jaeger: tracer.Jaeger{
			Enabled:      true,
			CollectorURL: "http://localhost:14268/api/traces",
		},
		SigNoz: tracer.SigNoz{
			Enabled:      false,
			CollectorURL: "192.168.1.5:4317",
		},
	}, "Service Auth")

	if trc.Err != nil {
		panic(trc.Err)
	}

	//Nats
	natsServer := nats.New(nats.Options{
		Host: "localhost",
		Port: "4222",
	})

	if natsServer.Error != nil {
		panic(natsServer.Error)
	}

	events.Setup(natsServer.Conn)

	fmt.Println("RUN AUTH MICRO-SERVICE")
	runtime.Goexit()
}
