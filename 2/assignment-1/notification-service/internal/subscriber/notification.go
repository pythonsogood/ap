package subscriber

import (
	"log"

	"github.com/nats-io/nats.go"
)

type EventSubscriber interface {
	SubscribeNotification(subj string) error
}

type natsEventSubscriber struct {
	nc *nats.Conn
}

func NewNATSEventSubscriber(nc *nats.Conn) EventSubscriber {
	return &natsEventSubscriber{
		nc: nc,
	}
}

func (e *natsEventSubscriber) SubscribeNotification(subj string) error {
	_, err := e.nc.Subscribe(subj, func(m *nats.Msg) {
		log.Printf("%v\n", string(m.Data))
	})

	if err != nil {
		return err
	}

	return nil
}
