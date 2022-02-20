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
		AgentHost:    "localhost",
		AgentPort:    "6831",
		AgentConnect: true,
	}, "Service Auth")
	defer trc.Finish()

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
