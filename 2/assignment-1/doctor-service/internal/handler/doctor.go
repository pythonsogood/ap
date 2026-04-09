package handler

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/pythonsogood/ap-assignment1/proto"
)

type DoctorHTTPHandler interface {
	GETByID(c *gin.Context)
	GETList(c *gin.Context)
	POST(c *gin.Context)
}

type doctorHTTPHandlerImpl struct {
	service service.DoctorService
}

func NewDoctorHTTPHandler(service service.DoctorService) DoctorHTTPHandler {
	return &doctorHTTPHandlerImpl{
		service: service,
	}
}

func (h *doctorHTTPHandlerImpl) GETByID(c *gin.Context) {
	id := c.Param("id")

	doctor, err := h.service.GetDoctor(id)

	if err != nil {
		status_code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			status_code = http.StatusNotFound
		}

		c.JSON(status_code, gin.H{"error": err.Error()})
		return
	}

	if doctor == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *doctorHTTPHandlerImpl) GETList(c *gin.Context) {
	doctors, err := h.service.GetAllDoctors()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

type DoctorHTTPPOSTBind struct {
	FullName       string `json:"full_name" binding:"required"`
	Specialization string `json:"specialization" binding:""`
	Email          string `json:"email" binding:"required,email"`
}

func (h *doctorHTTPHandlerImpl) POST(c *gin.Context) {
	var doctor_bind DoctorHTTPPOSTBind

	if err := c.ShouldBindJSON(&doctor_bind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.service.CreateDoctor(doctor_bind.FullName, doctor_bind.Specialization, doctor_bind.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

type DoctorGRPCHandler struct {
	pb.UnimplementedDoctorServiceServer
	service service.DoctorService
}

func NewDoctorGRPCHandler(service service.DoctorService) *DoctorGRPCHandler {
	return &DoctorGRPCHandler{
		service: service,
	}
}

func (h *DoctorGRPCHandler) GetDoctor(_ context.Context, in *pb.GetDoctorRequest) (*pb.DoctorResponse, error) {
	id := in.GetId()

	doctor, err := h.service.GetDoctor(id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Doctor with id %s not found", id)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}, nil
}

func (h *DoctorGRPCHandler) GetDoctors(_ context.Context, in *pb.GetDoctorsRequest) (*pb.GetDoctorsResponse, error) {
	doctors, err := h.service.GetAllDoctors()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	doctors_response := make([]*pb.DoctorResponse, 0)

	for _, doctor := range doctors {
		doctors_response = append(doctors_response, &pb.DoctorResponse{
			Id:             doctor.ID,
			FullName:       doctor.FullName,
			Specialization: doctor.Specialization,
			Email:          doctor.Email,
		})
	}

	return &pb.GetDoctorsResponse{
		Doctors: doctors_response,
	}, nil
}

func (h *DoctorGRPCHandler) CreateDoctor(_ context.Context, in *pb.CreateDoctorRequest) (*pb.DoctorResponse, error) {
	doctor, err := h.service.CreateDoctor(in.GetFullName(), in.GetSpecialization(), in.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}, nil
}
