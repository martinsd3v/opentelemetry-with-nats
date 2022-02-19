package nats

import (
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"

	"github.com/nats-io/nats.go"
)

//RespondDto data transfer object
type RespondDto struct {
	Data        interface{}
	SpanContext *tracer.SpanContext
	NatsMsg     nats.Msg
}

//Respond nats message respond
func Respond(dto RespondDto) error {
	return dto.NatsMsg.Respond(DataToByte(dto.SpanContext, dto.Data))
}
