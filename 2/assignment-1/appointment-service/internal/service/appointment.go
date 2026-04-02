package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
)

type AppointmentService interface {
	GetAppointment(id string) (*model.Appointment, error)
	GetAllAppointments() ([]*model.Appointment, error)
	CreateAppointment(full_name string, specialization string, email string) (*model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
	isValidDoctorId(doctor_id string) (bool, error)
}

type appointmentServiceImpl struct {
	repo                   repository.AppointmentRepository
	doctor_service_address string
	HTTPClient             *http.Client
}

func (s appointmentServiceImpl) isValidDoctorId(doctor_id string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/doctors/%s", s.doctor_service_address, doctor_id), nil)

	if err != nil {
		return false, err
	}

	req.Header.Set("X-Service-Name", "appointment-service")

	resp, err := s.HTTPClient.Do(req)

	if err != nil {
		return false, fmt.Errorf("Doctors service is currently unavailable\nError: %s", err.Error())
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s appointmentServiceImpl) GetAppointment(id string) (*model.Appointment, error) {
	return s.repo.FindByID(id)
}

func (s appointmentServiceImpl) GetAllAppointments() ([]*model.Appointment, error) {
	return s.repo.All()
}

func (s appointmentServiceImpl) CreateAppointment(title string, description string, doctor_id string) (*model.Appointment, error) {
	valid_doctor_id, err := s.isValidDoctorId(doctor_id)

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

	if appointment.Status == model.StatusDone {
		return errors.New("Done appointments cannot be updated")
	}

	return s.repo.UpdateStatus(id, status)
}

func NewAppointmentService(repo repository.AppointmentRepository, doctor_service_address string, http_client *http.Client) AppointmentService {
	if http_client == nil {
		http_client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	return &appointmentServiceImpl{
		repo:                   repo,
		doctor_service_address: doctor_service_address,
		HTTPClient:             http_client,
	}
}
