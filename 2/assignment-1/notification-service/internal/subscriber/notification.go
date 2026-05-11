package subscriber

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/notification/internal/jobqueue"
)

type EventSubscriber interface {
	SubscribeNotification(subj string) error
}

type natsEventSubscriber struct {
	nc    *nats.Conn
	queue *jobqueue.Queue
}

func NewNATSEventSubscriber(nc *nats.Conn, queue *jobqueue.Queue) EventSubscriber {
	return &natsEventSubscriber{
		nc:    nc,
		queue: queue,
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

		if subj == "appointments.status_updated" {
			newStatus, _ := data["new_status"].(string)

			if newStatus == "DONE" {
				eventType, _ := data["event_type"].(string)
				appointmentID, _ := data["id"].(string)
				occurredAt, _ := data["occurred_at"].(string)
				doctorID, _ := data["doctor_id"].(string)

				idempotencyKey := jobqueue.BuildIdempotencyKey(eventType, appointmentID, occurredAt)
				job := jobqueue.BuildJob(appointmentID, doctorID, occurredAt, idempotencyKey)

				e.queue.Enqueue(context.Background(), job)
			}
		}
	})

	if err != nil {
		return err
	}

	return nil
}
