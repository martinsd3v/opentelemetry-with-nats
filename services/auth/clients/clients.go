package clients

import (
	"context"

	"opentelemetry/services/auth/events"
	natsUtil "opentelemetry/utils/nats"
	"opentelemetry/utils/tracer"

	"github.com/nats-io/nats.go"
)

type client struct {
	conn *nats.Conn
}

func Setup(conn *nats.Conn) client {
	return client{conn: conn}
}

func (c *client) Auth(ctx context.Context, request events.AuthRequest) (events.AuthResponse, error) {
	ctx, span := tracer.New(ctx).WithNewTrace("ServiceAuth", "clients/Auth")
	defer span.Finish()

	spanContext := span.ExportSpanContext()
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
