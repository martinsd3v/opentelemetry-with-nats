package main

import (
	"fmt"
	"runtime"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
)

func main() {
	//Tracing
	tracing := tracer.SetupJeagerTracer(tracer.Options{
		EndpointURL: "http://localhost:14268/api/traces",
	})
	if tracing.Error != nil {
		panic(tracing.Error)
	}

	//Nats
	natsServer := nats.New(nats.Options{
		Host: "localhost",
		Port: "4222",
	})

	if natsServer.Error != nil {
		panic(natsServer.Error)
	}

	events.Setup(natsServer.Conn, tracing)

	fmt.Println("RUN AUTH MICRO-SERVICE")
	runtime.Goexit()
}
