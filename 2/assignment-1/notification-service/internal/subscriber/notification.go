package subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

		message := fmt.Sprintf("{\"time\": \"%s\",\"subject\":\"%s\",\"event\":{", time.Now().Format(time.RFC3339), subj)

		for k, v := range data {
			message += fmt.Sprintf("\"%s\": %v,", k, v)
			log.Printf("%s: %v\n", k, v)
		}

		message = strings.TrimSuffix(message, ",") + "}}"

		log.Printf("%v\n", message)
	})

	if err != nil {
		return err
	}

	return nil
}
