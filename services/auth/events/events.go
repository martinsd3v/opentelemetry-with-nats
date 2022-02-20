package events

import (
	"context"

	useCase "github.com/martinsd3v/opentelemetry-with-nats/services/auth/use_cases"
	natsUtil "github.com/martinsd3v/opentelemetry-with-nats/utils/nats"
	"github.com/martinsd3v/opentelemetry-with-nats/utils/open_telemetry/provider"

	"github.com/nats-io/nats.go"
)

type event struct {
	conn *nats.Conn
}

const (
	QueueAuth = "queue-auth"
)

func Setup(conn *nats.Conn) {
	e := event{conn}
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

		ctx := context.Background()
		ctx = provider.InjectContext(ctx, spanConfig)
		ctx, span := provider.Span(ctx, "events/auth", provider.SpanStartOption{
			Key:   "Receive Data From Nats",
			Value: request,
		})
		defer span.End()

		dto := natsUtil.RespondDto{
			SpanContext: spanConfig,
			NatsMsg:     *msg,
			Data:        response,
		}

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
