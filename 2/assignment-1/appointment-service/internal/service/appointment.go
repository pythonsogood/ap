package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
)

type AppointmentService interface {
	GetAppointment(id string) (*model.Appointment, error)
	GetAllAppointments() ([]*model.Appointment, error)
	CreateAppointment(full_name string, specialization string, email string) (*model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
}

type appointmentServiceImpl struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentServiceImpl{
		repo: repo,
	}
}

func (s appointmentServiceImpl) GetAppointment(id string) (*model.Appointment, error) {
	return s.repo.FindByID(id)
}

func (s appointmentServiceImpl) GetAllAppointments() ([]*model.Appointment, error) {
	return s.repo.All()
}

func (s appointmentServiceImpl) CreateAppointment(title string, description string, doctor_id string) (*model.Appointment, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	id := uuid.String()

	appointment := model.NewAppointment(id, title, description, doctor_id)

	if err := s.repo.Insert(appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (s appointmentServiceImpl) UpdateStatus(id string, status model.Status) error {
	appointment, err := s.repo.FindByID(id)

	if err != nil {
		return err
	}

	if appointment.Status == model.StatusDone {
		return errors.New("Done appointments cannot be updated")
	}

	return s.repo.UpdateStatus(id, status)
}
