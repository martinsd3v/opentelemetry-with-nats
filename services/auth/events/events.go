package events

import (
	"context"

	useCase "github.com/martinsd3v/opentelemetry-with-nats/services/auth/use_cases"
	natsUtil "github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/tracer"

	"github.com/nats-io/nats.go"
)

type event struct {
	conn    *nats.Conn
	tracing tracer.Tracing
}

const (
	QueueAuth = "queue-auth"
)

func Setup(conn *nats.Conn, tracing tracer.Tracing) {
	e := event{conn, tracing}
	conn.QueueSubscribe(QueueAuth, "queue", e.auth)
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Auth  bool  `json:"auth"`
	Error error `json:"error"`
}

func (e *event) auth(msg *nats.Msg) {
	if msg != nil {
		request, response := AuthRequest{}, AuthResponse{}
		spanConfig, err := natsUtil.ByteToData(msg.Data, &request)

		ctx := tracer.ContextFromSpanContext(context.Background(), spanConfig)
		ctx, span := e.tracing.New(ctx).WithNewTrace("ServiceAuth", "events/auth", tracer.SpanStartOption{
			Key:   "Receive data from nats",
			Value: request,
		})
		defer span.Finish()

		dto := natsUtil.RespondDto{
			SpanContext: spanConfig,
			NatsMsg:     *msg,
			Data:        response,
		}
		natsUtil.Respond(dto)

		if err != nil {
			response.Error = err
			dto.Data = response
			natsUtil.Respond(dto)
			return
		}

		services := useCase.New(e.tracing)
		services.HashPassword(ctx)
		response.Auth = services.AuthUser(ctx, request.Email, request.Password)

		dto.Data = response
		natsUtil.Respond(dto)
	}
}
