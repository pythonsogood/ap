package service

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type DoctorService interface {
	IsValidDoctorId(doctor_id string) (bool, error)
}

type doctorServiceImpl struct {
	doctor_service_address string
	HTTPClient             *http.Client
}

func (s doctorServiceImpl) IsValidDoctorId(doctor_id string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/doctors/%s", s.doctor_service_address, doctor_id), nil)

	if err != nil {
		return false, err
	}

	req.Header.Set("X-Service-Name", "appointment-service")

	resp, err := s.HTTPClient.Do(req)

	if err != nil {
		log.Println("[ERROR] Doctors service is currently unavailable\nError: ", err.Error())

		return false, fmt.Errorf("Doctors service is currently unavailable\nError: %s", err.Error())
	}

	return resp.StatusCode == http.StatusOK, nil
}

func NewDoctorService(doctor_service_address string, http_client *http.Client) DoctorService {
	if http_client == nil {
		http_client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	return &doctorServiceImpl{
		doctor_service_address: doctor_service_address,
		HTTPClient:             http_client,
	}
}
