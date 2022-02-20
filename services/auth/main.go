package main

import (
	"fmt"
	"runtime"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"
)

func main() {
	//Tracer
	trc := provider.Start(provider.Options{
		EndpointURL: "http://localhost:14268/api/traces",
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

	events.Setup(natsServer.Conn, trc)

	fmt.Println("RUN AUTH MICRO-SERVICE")
	runtime.Goexit()
}
