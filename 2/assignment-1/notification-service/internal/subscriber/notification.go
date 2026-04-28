package subscriber

import "github.com/nats-io/nats.go"

type EventSubscriber interface{}

type natsEventSubscriber struct{}

func NewNATSEventSubscriber(nc *nats.Conn) EventSubscriber {
	return &natsEventSubscriber{}
}
