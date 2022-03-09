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
		EndpointURL: "localhost:14250",
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
