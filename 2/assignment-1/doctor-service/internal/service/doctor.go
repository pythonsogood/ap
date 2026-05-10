package service

import (
	"github.com/google/uuid"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/cache"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/repository"
)

type DoctorService interface {
	GetDoctor(id string) (*model.Doctor, error)
	GetAllDoctors() ([]*model.Doctor, error)
	CreateDoctor(full_name string, specialization string, email string) (*model.Doctor, error)
}

type doctorServiceImpl struct {
	repo  repository.DoctorRepository
	cache cache.DoctorCacheRepository
}

func NewDoctorService(repo repository.DoctorRepository, cache cache.DoctorCacheRepository) DoctorService {
	return &doctorServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

func (s doctorServiceImpl) GetDoctor(id string) (*model.Doctor, error) {
	cached, err := s.cache.GetDoctor(id)

	if err == nil && cached != nil {
		return cached, nil
	}

	return s.repo.FindByID(id)
}

func (s doctorServiceImpl) GetAllDoctors() ([]*model.Doctor, error) {
	return s.repo.All()
}

func (s doctorServiceImpl) CreateDoctor(full_name string, specialization string, email string) (*model.Doctor, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	id := uuid.String()

	doctor := model.NewDoctor(id, full_name, specialization, email)

	if err := s.repo.Insert(doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}
