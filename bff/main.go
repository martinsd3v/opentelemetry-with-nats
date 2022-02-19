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
	//Tracing
	tracing := tracer.SetupJeagerTracer(tracer.Options{
		EndpointURL: "http://localhost:14268/api/traces",
	})
	if tracing.Error != nil {
		panic(tracing.Error)
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
		ctx, span := tracing.New(ctx).WithNewTrace("Back For Front", "route/auth")
		defer span.Finish()

		_, sp2 := tracing.New(ctx).Simple("span2")
		defer sp2.Finish()

		_, sp44 := tracing.New(ctx).Simple("span2")
		defer sp44.Finish()

		ka2, spa33s := tracing.New(ctx).Trace("Nome aqui").Simple("Hahahaha")
		defer spa33s.Finish()

		_, sp443 := tracing.New(ka2).Simple("hehehehe")
		defer sp443.Finish()

		ctx2, sp1 := tracing.New(ctx).WithNewTrace("Trace1", "span1")
		defer sp1.Finish()

		ctx, sp3 := tracing.New(ctx2).WithNewTrace("Trace3", "span3")
		defer sp3.Finish()

		response, err := natsClients.Auth(ctx, tracing, events.AuthRequest{
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
