package main

import (
	"fmt"
	"opentelemetry/services/auth/events"
	"opentelemetry/utils/nats"
	"runtime"
)

func main() {
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
