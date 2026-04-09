package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/pythonsogood/ap-assignment1/proto"
)

type DoctorService interface {
	IsValidDoctorId(doctor_id string) (bool, error)
}

type httpDoctorServiceImpl struct {
	doctor_service_address string
	HTTPClient             *http.Client
}

func (s httpDoctorServiceImpl) IsValidDoctorId(doctor_id string) (bool, error) {
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

func NewHTTPDoctorService(doctor_service_address string, http_client *http.Client) DoctorService {
	if http_client == nil {
		http_client = &http.Client{
			Timeout: 15 * time.Second,
		}
	}

	return &httpDoctorServiceImpl{
		doctor_service_address: doctor_service_address,
		HTTPClient:             http_client,
	}
}

type grpcDoctorServiceImpl struct {
	doctor_service_address string
	timeout                time.Duration
}

func (s grpcDoctorServiceImpl) IsValidDoctorId(doctor_id string) (bool, error) {
	conn, err := grpc.NewClient(s.doctor_service_address, grpc.WithTransportCredentials(insecure.NewBundle().TransportCredentials()))

	if err != nil {
		return false, err
	}

	defer conn.Close()

	client := pb.NewDoctorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err = client.GetDoctor(ctx, &pb.GetDoctorRequest{
		Id: doctor_id,
	})

	if err != nil {
		if strings.Contains(err.Error(), "code = NotFound") {
			return false, nil
		}

		return false, fmt.Errorf("[ERROR] Doctors service is currently unavailable\nError: %s", err.Error())
	}

	return true, nil
}

func NewGRPCDoctorService(doctor_service_address string, timeout time.Duration) DoctorService {
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	return &grpcDoctorServiceImpl{
		doctor_service_address: doctor_service_address,
		timeout:                timeout,
	}
}
