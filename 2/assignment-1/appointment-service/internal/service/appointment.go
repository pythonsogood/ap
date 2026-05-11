package service

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/cache"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
)

type AppointmentService interface {
	GetAppointment(id string) (*model.Appointment, error)
	GetAllAppointments() ([]*model.Appointment, error)
	CreateAppointment(title string, description string, doctor_id string) (*model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
}

type appointmentServiceImpl struct {
	repo           repository.AppointmentRepository
	cache          cache.AppointmentCacheRepository
	doctor_service DoctorService
}

func NewAppointmentService(repo repository.AppointmentRepository, cache cache.AppointmentCacheRepository, doctor_service DoctorService) AppointmentService {
	return &appointmentServiceImpl{
		repo:           repo,
		cache:          cache,
		doctor_service: doctor_service,
	}
}

func (s appointmentServiceImpl) GetAppointment(id string) (*model.Appointment, error) {
	if s.cache != nil {
		cached, found, err := s.cache.GetAppointment(id)

		if err != nil {
			log.Println(err.Error())
		} else {
			if found && cached != nil {
				return cached, nil
			}
		}
	}

	appointment, err := s.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.cache.SaveAppointment(appointment); err != nil {
			log.Println(err.Error())
		}
	}

	return appointment, nil
}

func (s appointmentServiceImpl) GetAllAppointments() ([]*model.Appointment, error) {
	if s.cache != nil {
		cached, found, err := s.cache.GetAppointments()

		if err != nil {
			log.Println(err.Error())
		} else {
			if found && cached != nil {
				return cached, nil
			}
		}
	}

	appointments, err := s.repo.All()

	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.cache.SaveAppointments(appointments); err != nil {
			log.Println(err.Error())
		}
	}

	return appointments, nil
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

	if s.cache != nil {
		if err := s.cache.CreateAppointment(appointment); err != nil {
			log.Println(err.Error())
		}
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

	err = s.repo.UpdateStatus(id, status)

	if err != nil {
		return err
	}

	if s.cache != nil {
		appointment, err = s.repo.FindByID(appointment.ID)

		if err == nil {
			if err := s.cache.UpdateAppointmentStatus(appointment); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}

	return nil
}
