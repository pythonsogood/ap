package event

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

type EventPublisher interface {
	DoctorCreated(event *DoctorCreatedEvent) error
}

type DoctorCreatedEvent struct {
	EventType      string `json:"event_type"`
	OccurredAt     string `json:"occurred_at"`
	DoctorId       string `json:"id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
}

func NewDoctorCreatedEvent(doctor_id string, full_name string, specialization string, email string, occurred_at time.Time) *DoctorCreatedEvent {
	if occurred_at.IsZero() {
		occurred_at = time.Now()
	}

	return &DoctorCreatedEvent{
		EventType:      "doctors.created",
		OccurredAt:     occurred_at.String(),
		DoctorId:       doctor_id,
		FullName:       full_name,
		Specialization: specialization,
		Email:          email,
	}
}

type NATSEventPublisher struct {
	nc *nats.Conn
}

func NewNATSEventPublisher(nc *nats.Conn) *NATSEventPublisher {
	return &NATSEventPublisher{
		nc: nc,
	}
}

func (e *NATSEventPublisher) DoctorCreated(event *DoctorCreatedEvent) error {
	event_json, err := json.Marshal(event)

	if err != nil {
		return err
	}

	if err := e.nc.Publish("doctors.created", event_json); err != nil {
		return err
	}

	return nil
}
