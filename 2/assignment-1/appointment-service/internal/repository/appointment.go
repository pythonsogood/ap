package repository

import (
	"database/sql"
	"time"

	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
)

type AppointmentRepository interface {
	FindByID(ID string) (*model.Appointment, error)
	All() ([]*model.Appointment, error)
	Insert(appointment *model.Appointment) error
	UpdateStatus(ID string, status model.Status) error
}

type SQLiteAppointmentRepository struct {
	db *sql.DB
}

func NewSQLiteAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &SQLiteAppointmentRepository{
		db: db,
	}
}

func (r *SQLiteAppointmentRepository) FindByID(ID string) (*model.Appointment, error) {
	statement, err := r.db.Prepare("SELECT id, title, description, doctor_id, status, created_at, updated_at FROM appointments WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var id, title, description, doctor_id string
	var status model.Status
	var created_at, updated_at time.Time

	if err := statement.QueryRow(ID).Scan(&id, &title, &description, &doctor_id, &status, &created_at, &updated_at); err != nil {
		return nil, err
	}

	return &model.Appointment{
		ID:          id,
		Title:       title,
		Description: description,
		DoctorID:    doctor_id,
		Status:      status,
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}, nil
}

func (r *SQLiteAppointmentRepository) All() ([]*model.Appointment, error) {
	rows, err := r.db.Query("SELECT id, title, description, doctor_id, status, created_at, updated_at FROM appointments")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	appointments := make([]*model.Appointment, 0)

	for rows.Next() {
		var id, title, description, doctor_id string
		var status model.Status
		var created_at, updated_at time.Time

		if err := rows.Scan(&id, &title, &description, &doctor_id, &status, &created_at, &updated_at); err == nil {
			appointment := &model.Appointment{
				ID:          id,
				Title:       title,
				Description: description,
				DoctorID:    doctor_id,
				Status:      status,
				CreatedAt:   created_at,
				UpdatedAt:   updated_at,
			}

			appointments = append(appointments, appointment)
		}
	}

	return appointments, nil
}

func (r *SQLiteAppointmentRepository) Insert(appointment *model.Appointment) error {
	statement, err := r.db.Prepare("INSERT INTO appointments(id, title, description, doctor_id, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(appointment.ID, appointment.Title, appointment.Description, appointment.DoctorID, appointment.Status, appointment.CreatedAt, appointment.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLiteAppointmentRepository) UpdateStatus(ID string, status model.Status) error {
	statement, err := r.db.Prepare("UPDATE appointments SET status = ?, updated_at = ? WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(status, time.Now(), ID)

	if err != nil {
		return err
	}

	return nil
}

type PostgresAppointmentRepository struct {
	db *sql.DB
}

func NewPostgresAppointmentRepository(db *sql.DB) AppointmentRepository {
	return &PostgresAppointmentRepository{
		db: db,
	}
}

func (r *PostgresAppointmentRepository) FindByID(ID string) (*model.Appointment, error) {
	statement, err := r.db.Prepare("SELECT id, title, description, doctor_id, status, created_at, updated_at FROM appointments WHERE id = $1")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	var id, title, description, doctor_id string
	var status model.Status
	var created_at, updated_at time.Time

	if err := statement.QueryRow(ID).Scan(&id, &title, &description, &doctor_id, &status, &created_at, &updated_at); err != nil {
		return nil, err
	}

	return &model.Appointment{
		ID:          id,
		Title:       title,
		Description: description,
		DoctorID:    doctor_id,
		Status:      status,
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}, nil
}

func (r *PostgresAppointmentRepository) All() ([]*model.Appointment, error) {
	rows, err := r.db.Query("SELECT id, title, description, doctor_id, status, created_at, updated_at FROM appointments")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	appointments := make([]*model.Appointment, 0)

	for rows.Next() {
		var id, title, description, doctor_id string
		var status model.Status
		var created_at, updated_at time.Time

		if err := rows.Scan(&id, &title, &description, &doctor_id, &status, &created_at, &updated_at); err == nil {
			appointment := &model.Appointment{
				ID:          id,
				Title:       title,
				Description: description,
				DoctorID:    doctor_id,
				Status:      status,
				CreatedAt:   created_at,
				UpdatedAt:   updated_at,
			}

			appointments = append(appointments, appointment)
		}
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) Insert(appointment *model.Appointment) error {
	statement, err := r.db.Prepare("INSERT INTO appointments(id, title, description, doctor_id, status, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(appointment.ID, appointment.Title, appointment.Description, appointment.DoctorID, appointment.Status, appointment.CreatedAt, appointment.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresAppointmentRepository) UpdateStatus(ID string, status model.Status) error {
	statement, err := r.db.Prepare("UPDATE appointments SET status = $1, updated_at = $2 WHERE id = $3")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(status, time.Now(), ID)

	if err != nil {
		return err
	}

	return nil
}
