package subscriber

import (
	"encoding/json"
	"log"
	"time"

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
		var data map[string]any

		if err := json.Unmarshal(m.Data, &data); err != nil {
			log.Printf("[SubscribeNotification] json unmarshal error: %s. Raw data: %s\n", err.Error(), string(m.Data))
			return
		}

		message := map[string]any{
			"time":    time.Now().UTC().Format(time.RFC3339),
			"subject": subj,
			"event":   data,
		}

		messageJSON, err := json.Marshal(message)

		if err != nil {
			log.Printf("[SubscribeNotification] json marshal error: %s\n", err.Error())
			return
		}

		log.Println(string(messageJSON))
	})

	if err != nil {
		return err
	}

	return nil
}
