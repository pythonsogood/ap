package model

import (
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Appointment struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DoctorID    string    `json:"doctor_id"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewAppointment(id string, title string, description string, doctor_id string) *Appointment {
	return &Appointment{
		ID:          id,
		Title:       title,
		Description: description,
		DoctorID:    doctor_id,
		Status:      StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
