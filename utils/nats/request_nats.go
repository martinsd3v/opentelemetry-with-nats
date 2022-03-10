package nats

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"
	"github.com/nats-io/nats.go"
)

//RequestDto data transfer object
type RequestDto struct {
	Ctx         context.Context
	Queue       string
	Data        interface{}
	SpanContext *tracer.SpanContext
	NatsConn    *nats.Conn
}

//Request nats message request
func Request(dto RequestDto) (*nats.Msg, error) {
	data := DataToByte(dto.SpanContext, dto.Data)
	msg, err := dto.NatsConn.RequestWithContext(dto.Ctx, dto.Queue, data)

	return msg, err
}
