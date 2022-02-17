package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Options struct {
	Host string
	Port string
}

type natsConn struct {
	Conn  *nats.Conn
	Error error
}

var singletonNats *natsConn = nil

func New(ots Options) *natsConn {
	if singletonNats == nil {
		natsURL := fmt.Sprintf("nats://%s:%s", ots.Host, ots.Port)
		conn, err := nats.Connect(natsURL)
		return &natsConn{
			Conn:  conn,
			Error: err,
		}
	}
	return singletonNats
}
