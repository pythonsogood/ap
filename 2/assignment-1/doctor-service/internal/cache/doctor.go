package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	doctorKeyFormat = "doctor:%s"
	doctorsListKey  = "doctors:list"
)

type DoctorCacheRepository interface {
	GetDoctor(id string) (*model.Doctor, bool, error)
	GetDoctors() ([]*model.Doctor, bool, error)

	CreateDoctor() error
}

type redisDoctorCacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisDoctorCacheRepository(client *redis.Client, ttl time.Duration) DoctorCacheRepository {
	return &redisDoctorCacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisDoctorCacheRepository) GetDoctor(id string) (*model.Doctor, bool, error) {
	raw, err := r.client.Get(context.Background(), fmt.Sprintf(doctorKeyFormat, id)).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var doctor model.Doctor

	if err := json.Unmarshal(raw, &doctor); err != nil {
		return nil, false, err
	}

	return &doctor, true, nil
}

func (r *redisDoctorCacheRepository) GetDoctors() ([]*model.Doctor, bool, error) {
	raw, err := r.client.Get(context.Background(), doctorsListKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var doctors []*model.Doctor

	if err := json.Unmarshal(raw, &doctors); err != nil {
		return nil, false, err
	}

	return doctors, true, nil
}

func (r *redisDoctorCacheRepository) CreateDoctor() error {
	cmd := r.client.Del(context.Background(), doctorsListKey)

	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
