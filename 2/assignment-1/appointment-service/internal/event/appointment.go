package event

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
)

type EventPublisher interface {
	AppointmentCreated(event *AppointmentCreatedEvent) error
	AppointmentStatusUpdated(event *AppointmentStatusUpdatedEvent) error
}

type AppointmentCreatedEvent struct {
	EventType  string `json:"event_type"`
	OccurredAt string `json:"occurred_at"`
	Id         string `json:"id"`
	Title      string `json:"title"`
	DoctorId   string `json:"doctor_id"`
	Status     string `json:"status"`
}

func NewAppointmentCreatedEvent(id string, title string, doctor_id string, occurred_at time.Time) *AppointmentCreatedEvent {
	if occurred_at.IsZero() {
		occurred_at = time.Now()
	}

	return &AppointmentCreatedEvent{
		EventType:  "appointments.created",
		OccurredAt: occurred_at.Format(time.RFC3339),
		Id:         id,
		Title:      title,
		DoctorId:   doctor_id,
		Status:     string(model.StatusNew),
	}
}

type AppointmentStatusUpdatedEvent struct {
	EventType  string       `json:"event_type"`
	OccurredAt string       `json:"occurred_at"`
	Id         string       `json:"id"`
	OldStatus  model.Status `json:"old_status"`
	NewStatus  model.Status `json:"new_status"`
}

func NewAppointmentStatusUpdatedEvent(id string, old_status model.Status, new_status model.Status, occurred_at time.Time) *AppointmentStatusUpdatedEvent {
	if occurred_at.IsZero() {
		occurred_at = time.Now()
	}

	return &AppointmentStatusUpdatedEvent{
		EventType:  "appointments.status_updated",
		OccurredAt: occurred_at.Format(time.RFC3339),
		Id:         id,
		OldStatus:  old_status,
		NewStatus:  new_status,
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

func (e *NATSEventPublisher) AppointmentCreated(event *AppointmentCreatedEvent) error {
	event_json, err := json.Marshal(event)

	if err != nil {
		return err
	}

	if err := e.nc.Publish("appointments.created", event_json); err != nil {
		return err
	}

	return nil
}

func (e *NATSEventPublisher) AppointmentStatusUpdated(event *AppointmentStatusUpdatedEvent) error {
	event_json, err := json.Marshal(event)

	if err != nil {
		return err
	}

	if err := e.nc.Publish("appointments.status_updated", event_json); err != nil {
		return err
	}

	return nil
}
