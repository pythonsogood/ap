package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	doctorKeyFormat = "doctor:%s"
	doctorsListKey  = "doctors:list"
)

type DoctorCacheRepository interface {
	GetDoctor(id string) (*model.Doctor, error)
	GetDoctors() ([]*model.Doctor, error)
}

type redisDoctorCacheRepository struct {
	client *redis.Client
}

func NewRedisDoctorCacheRepository(client *redis.Client) DoctorCacheRepository {
	return &redisDoctorCacheRepository{
		client: client,
	}
}

func (r *redisDoctorCacheRepository) GetDoctor(id string) (*model.Doctor, error) {
	raw, err := r.client.Get(context.Background(), fmt.Sprintf(doctorKeyFormat, id)).Bytes()

	if err != nil {
		return nil, err
	}

	var doctor model.Doctor

	if err := json.Unmarshal(raw, &doctor); err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *redisDoctorCacheRepository) GetDoctors() ([]*model.Doctor, error) {
	raw, err := r.client.Get(context.Background(), doctorsListKey).Bytes()

	if err != nil {
		return nil, err
	}

	var doctors []*model.Doctor

	if err := json.Unmarshal(raw, &doctors); err != nil {
		return nil, err
	}

	return doctors, nil
}
