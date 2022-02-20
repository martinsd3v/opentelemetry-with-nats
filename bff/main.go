package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/clients"
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
	}, "Back For Front")
	defer trc.Finish()

	if trc.Err != nil {
		panic(trc.Err)
	}

	//Nats
	natsServer := nats.New(nats.Options{Host: "localhost", Port: "4222"})

	if natsServer.Error != nil {
		panic(natsServer.Error)
	}

	natsClients := clients.Setup(natsServer.Conn)

	//Echo
	e := echo.New()
	e.POST("/auth", func(c echo.Context) error {
		ctx := c.Request().Context()

		ctx, span := provider.Span(ctx, "route/auth")
		defer span.End()

		response, err := natsClients.Auth(ctx, events.AuthRequest{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, response)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
