package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	appointmentKeyFormat = "appointment:%s"
	appointmentsListKey  = "appointments:list"
)

type AppointmentCacheRepository interface {
	GetAppointment(id string) (*model.Appointment, bool, error)
	SaveAppointment(appointment *model.Appointment) error

	GetAppointments() ([]*model.Appointment, bool, error)
	SaveAppointments(appointments []*model.Appointment) error

	CreateAppointment(appointment *model.Appointment) error
	UpdateAppointmentStatus(appointment *model.Appointment) error
}

type redisAppointmentCacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisAppointmentCacheRepository(client *redis.Client, ttl time.Duration) AppointmentCacheRepository {
	return &redisAppointmentCacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisAppointmentCacheRepository) GetAppointment(id string) (*model.Appointment, bool, error) {
	raw, err := r.client.Get(context.Background(), fmt.Sprintf(appointmentKeyFormat, id)).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var appointment model.Appointment

	if err := json.Unmarshal(raw, &appointment); err != nil {
		return nil, false, err
	}

	return &appointment, true, nil
}

func (r *redisAppointmentCacheRepository) SaveAppointment(appointment *model.Appointment) error {
	raw, err := json.Marshal(appointment)

	if err != nil {
		return err
	}

	cmd := r.client.Set(context.Background(), fmt.Sprintf(appointmentKeyFormat, appointment.ID), raw, r.ttl)

	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *redisAppointmentCacheRepository) GetAppointments() ([]*model.Appointment, bool, error) {
	raw, err := r.client.Get(context.Background(), appointmentsListKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	var appointments []*model.Appointment

	if err := json.Unmarshal(raw, &appointments); err != nil {
		return nil, false, err
	}

	return appointments, true, nil
}

func (r *redisAppointmentCacheRepository) SaveAppointments(appointments []*model.Appointment) error {
	raw, err := json.Marshal(appointments)

	if err != nil {
		return err
	}

	cmd := r.client.Set(context.Background(), appointmentsListKey, raw, r.ttl)

	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *redisAppointmentCacheRepository) CreateAppointment(appointment *model.Appointment) error {
	cmd := r.client.Del(context.Background(), appointmentsListKey)

	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *redisAppointmentCacheRepository) UpdateAppointmentStatus(appointment *model.Appointment) error {
	if err := r.SaveAppointment(appointment); err != nil {
		return err
	}

	cmd := r.client.Del(context.Background(), appointmentsListKey)

	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
