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

const AppointmentTableName = "appointments"
const AppointmentTableCreateQuery = `
	CREATE TABLE IF NOT EXISTS ` + AppointmentTableName + ` (
		id TEXT UNIQUE NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		doctor_id TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)
`

type Appointment struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DoctorID    string    `json:"doctor_id"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Appointment) TableName() string {
	return AppointmentTableName
}

func (a *Appointment) TableCreateQuery() string {
	return AppointmentTableCreateQuery
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
