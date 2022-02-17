package nats

import (
	"context"

	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"

	"github.com/nats-io/nats.go"
)

//RespondDto data transfer object
type RespondDto struct {
	Ctx         context.Context
	Data        interface{}
	SpanContext *tracer.SpanContext
	NatsMsg     nats.Msg
}

//Respond nats message respond
func Respond(dto RespondDto) {
	if err := dto.NatsMsg.Respond(DataToByte(dto.SpanContext, dto.Data)); err != nil {
		_, span := tracer.New(dto.Ctx).WithNewTrace("Nats/Respond", "nats/Error")
		defer span.Finish()
	}
}
