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
	repo           repository.AppointmentRepository
	doctor_service DoctorService
}

func (s appointmentServiceImpl) GetAppointment(id string) (*model.Appointment, error) {
	return s.repo.FindByID(id)
}

func (s appointmentServiceImpl) GetAllAppointments() ([]*model.Appointment, error) {
	return s.repo.All()
}

func (s appointmentServiceImpl) CreateAppointment(title string, description string, doctor_id string) (*model.Appointment, error) {
	valid_doctor_id, err := s.doctor_service.IsValidDoctorId(doctor_id)

	if err != nil {
		return nil, err
	}

	if !valid_doctor_id {
		return nil, errors.New("Invalid doctor id")
	}

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

	valid_doctor_id, err := s.doctor_service.IsValidDoctorId(appointment.DoctorID)

	if err != nil {
		return err
	}

	if !valid_doctor_id {
		return errors.New("Invalid doctor id")
	}

	if appointment.Status == model.StatusDone {
		return errors.New("Done appointments cannot be updated")
	}

	return s.repo.UpdateStatus(id, status)
}

func NewAppointmentService(repo repository.AppointmentRepository, doctor_service DoctorService) AppointmentService {
	return &appointmentServiceImpl{
		repo:           repo,
		doctor_service: doctor_service,
	}
}
