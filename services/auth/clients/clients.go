package clients

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/services/auth/events"
	natsUtil "github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"

	"github.com/nats-io/nats.go"
)

type client struct {
	conn *nats.Conn
}

func Setup(conn *nats.Conn) client {
	return client{conn: conn}
}

func (c *client) Auth(ctx context.Context, tracing tracer.Tracing, request events.AuthRequest) (events.AuthResponse, error) {
	ctx, span := tracing.New(ctx).WithNewTrace("ServiceAuth", "clients/Auth")
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
