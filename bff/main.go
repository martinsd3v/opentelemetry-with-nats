package main

import (
	"net/http"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/clients"
	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"

	"github.com/labstack/echo/v4"
)

func main() {
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
		ctx, span := tracer.New(ctx).WithNewTrace("Back For Front", "route/auth")
		defer span.Finish()

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
