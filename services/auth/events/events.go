package events

import (
	"context"

	useCase "opentelemetry/services/auth/use_cases"
	natsUtil "opentelemetry/utils/nats"
	"opentelemetry/utils/tracer"

	"github.com/nats-io/nats.go"
)

type event struct {
	conn *nats.Conn
}

const (
	QueueAuth = "queue-auth"
)

func Setup(conn *nats.Conn) {
	e := event{conn: conn}
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
		ctx, span := tracer.New(ctx).WithNewTrace("ServiceAuth", "events/auth", tracer.SpanStartOption{
			Key:   "Receive data from nats",
			Value: request,
		})
		defer span.Finish()

		dto := natsUtil.RespondDto{
			Ctx:         ctx,
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

		services := useCase.New()
		services.HashPassword(ctx)
		response.Auth = services.AuthUser(ctx, request.Email, request.Password)

		dto.Data = response
		natsUtil.Respond(dto)
	}
}
