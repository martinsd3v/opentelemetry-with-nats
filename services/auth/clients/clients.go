package clients

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	natsUtil "github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"

	"github.com/nats-io/nats.go"
)

type client struct {
	conn *nats.Conn
}

func Setup(conn *nats.Conn) client {
	return client{conn}
}

func (c *client) Auth(ctx context.Context, request events.AuthRequest) (events.AuthResponse, error) {
	ctx, span := provider.Span(ctx, "clients/Auths")
	defer span.End()

	spanContext := provider.ExportSpanContext(ctx)
	dto := natsUtil.RequestDto{
		Ctx:         ctx,
		Queue:       events.QueueAuth,
		Data:        request,
		SpanContext: &spanContext,
		NatsConn:    c.conn,
	}
	msg, err := natsUtil.Request(dto)

	var reponse events.AuthResponse
	if err == nil {
		natsUtil.ByteToData(msg.Data, &reponse)
	}

	return reponse, nil
}
