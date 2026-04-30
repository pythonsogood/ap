package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/event"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/pythonsogood/ap-assignment1/proto"
)

type AppointmentHTTPHandler interface {
	GETByID(c *gin.Context)
	GETList(c *gin.Context)
	POST(c *gin.Context)
	PATCHStatusByID(c *gin.Context)

	StatusValidator(fl validator.FieldLevel) bool
}

type appointmentHTTPHandlerImpl struct {
	service         service.AppointmentService
	event_publisher event.EventPublisher
}

func NewAppointmentHTTPHandler(service service.AppointmentService, event_publisher event.EventPublisher) AppointmentHTTPHandler {
	return &appointmentHTTPHandlerImpl{
		service:         service,
		event_publisher: event_publisher,
	}
}

func (h *appointmentHTTPHandlerImpl) GETByID(c *gin.Context) {
	id := c.Param("id")

	appointment, err := h.service.GetAppointment(id)

	if err != nil {
		status_code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			status_code = http.StatusNotFound
		}

		c.JSON(status_code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func (h *appointmentHTTPHandlerImpl) GETList(c *gin.Context) {
	appointments, err := h.service.GetAllAppointments()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

type AppointmentPOSTBind struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:""`
	DoctorID    string `json:"doctor_id" binding:"required"`
}

func (h *appointmentHTTPHandlerImpl) POST(c *gin.Context) {
	var appointment_bind AppointmentPOSTBind

	if err := c.ShouldBindJSON(&appointment_bind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.service.CreateAppointment(appointment_bind.Title, appointment_bind.Description, appointment_bind.DoctorID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.event_publisher.AppointmentCreated(
		event.NewAppointmentCreatedEvent(
			appointment.ID,
			appointment.Title,
			appointment.DoctorID,
			time.Now(),
		),
	)

	if err != nil {
		log.Printf("[CreateDoctor] event publisher error: %s\n", err.Error())
	}

	c.JSON(http.StatusOK, appointment)
}

type AppointmentPATCHStatusBind struct {
	Status model.Status `json:"status" binding:"required,validstatus"`
}

func (h *appointmentHTTPHandlerImpl) PATCHStatusByID(c *gin.Context) {
	id := c.Param("id")

	var status_bind AppointmentPATCHStatusBind

	if err := c.ShouldBindJSON(&status_bind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	old_appointment, err := h.service.GetAppointment(id)

	if err != nil || old_appointment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment is not found"})
		return
	}

	old_appointment_status := old_appointment.Status
	new_appointment_status := status_bind.Status

	if err := h.service.UpdateStatus(id, new_appointment_status); err != nil {
		status_code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			status_code = http.StatusNotFound
		}

		c.JSON(status_code, gin.H{"error": err.Error()})
		return
	}

	err = h.event_publisher.AppointmentStatusUpdated(
		event.NewAppointmentStatusUpdatedEvent(
			id,
			old_appointment_status,
			new_appointment_status,
			time.Now(),
		),
	)

	if err != nil {
		log.Printf("[AppointmentStatusUpdated] event publisher error: %s\n", err.Error())
	}

	c.Status(http.StatusOK)
}

func (h *appointmentHTTPHandlerImpl) StatusValidator(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(model.Status)

	if !ok {
		return false
	}

	return status == model.StatusDone || status == model.StatusInProgress || status == model.StatusNew
}

type AppointmentGRPCHandler struct {
	pb.UnimplementedAppointmentServiceServer
	service         service.AppointmentService
	event_publisher event.EventPublisher
}

func NewAppointmentGRPCHandler(service service.AppointmentService, event_publisher event.EventPublisher) *AppointmentGRPCHandler {
	return &AppointmentGRPCHandler{
		service:         service,
		event_publisher: event_publisher,
	}
}

func (h *AppointmentGRPCHandler) GetAppointment(_ context.Context, in *pb.GetAppointmentRequest) (*pb.AppointmentResponse, error) {
	id := in.GetId()

	appointment, err := h.service.GetAppointment(id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Appointment with id %s not found", id)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	appointment_status, ok := pb.AppointmentStatus_value[strings.ToUpper(string(appointment.Status))]

	if !ok {
		appointment_status = int32(pb.AppointmentStatus_NEW)
	}

	return &pb.AppointmentResponse{
		Id:          appointment.ID,
		Title:       appointment.Title,
		Description: appointment.Description,
		DoctorId:    appointment.DoctorID,
		Status:      pb.AppointmentStatus(appointment_status),
		CreatedAt:   timestamppb.New(appointment.CreatedAt),
		UpdatedAt:   timestamppb.New(appointment.UpdatedAt),
	}, nil
}

func (h *AppointmentGRPCHandler) GetAppointments(_ context.Context, in *pb.GetAppointmentsRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := h.service.GetAllAppointments()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	appointments_response := make([]*pb.AppointmentResponse, len(appointments))

	for _, appointment := range appointments {
		appointment_status, ok := pb.AppointmentStatus_value[strings.ToUpper(string(appointment.Status))]

		if !ok {
			appointment_status = int32(pb.AppointmentStatus_NEW)
		}

		appointments_response = append(appointments_response, &pb.AppointmentResponse{
			Id:          appointment.ID,
			Title:       appointment.Title,
			Description: appointment.Description,
			DoctorId:    appointment.DoctorID,
			Status:      pb.AppointmentStatus(appointment_status),
			CreatedAt:   timestamppb.New(appointment.CreatedAt),
			UpdatedAt:   timestamppb.New(appointment.UpdatedAt),
		})
	}

	return &pb.GetAppointmentsResponse{
		Appointments: appointments_response,
	}, nil
}

func (h *AppointmentGRPCHandler) CreateAppointment(_ context.Context, in *pb.CreateAppointmentRequest) (*pb.AppointmentResponse, error) {
	appointment, err := h.service.CreateAppointment(in.GetTitle(), in.GetDescription(), in.GetDoctorId())

	if err != nil {
		code := codes.Internal

		if strings.Contains(err.Error(), "Doctors service is currently unavailable") {
			code = codes.Unavailable
		}

		if strings.Contains(err.Error(), "Invalid doctor id") {
			code = codes.FailedPrecondition
		}

		return nil, status.Error(code, err.Error())
	}

	appointment_status, ok := pb.AppointmentStatus_value[strings.ToUpper(string(appointment.Status))]

	if !ok {
		appointment_status = int32(pb.AppointmentStatus_NEW)
	}

	err = h.event_publisher.AppointmentCreated(
		event.NewAppointmentCreatedEvent(
			appointment.ID,
			appointment.Title,
			appointment.DoctorID,
			time.Now(),
		),
	)

	if err != nil {
		log.Printf("[AppointmentCreated] event publisher error: %s\n", err.Error())
	}

	return &pb.AppointmentResponse{
		Id:          appointment.ID,
		Title:       appointment.Title,
		Description: appointment.Description,
		DoctorId:    appointment.DoctorID,
		Status:      pb.AppointmentStatus(appointment_status),
		CreatedAt:   timestamppb.New(appointment.CreatedAt),
		UpdatedAt:   timestamppb.New(appointment.UpdatedAt),
	}, nil
}

func (h *AppointmentGRPCHandler) UpdateAppointmentStatus(_ context.Context, in *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error) {
	id := in.GetId()
	appointment_status, ok := pb.AppointmentStatus_name[int32(in.GetStatus())]

	if !ok {
		return nil, status.Error(codes.Unknown, "AppointmentStatus is not found")
	}

	old_appointment, err := h.service.GetAppointment(id)

	if err != nil || old_appointment == nil {
		return nil, status.Error(codes.Unknown, "Appointment is not found")
	}

	old_appointment_status := old_appointment.Status
	new_appointment_status := model.Status(appointment_status)

	if err := h.service.UpdateStatus(id, new_appointment_status); err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Appointment with id %s not found", id)
		}

		code := codes.Internal

		if strings.Contains(err.Error(), "Done appointments cannot be updated") {
			code = codes.InvalidArgument
		}

		if strings.Contains(err.Error(), "Doctors service is currently unavailable") {
			code = codes.Unavailable
		}

		if strings.Contains(err.Error(), "Invalid doctor id") {
			code = codes.FailedPrecondition
		}

		return nil, status.Error(code, err.Error())
	}

	err = h.event_publisher.AppointmentStatusUpdated(
		event.NewAppointmentStatusUpdatedEvent(
			id,
			old_appointment_status,
			new_appointment_status,
			time.Now(),
		),
	)

	if err != nil {
		log.Printf("[AppointmentStatusUpdated] event publisher error: %s\n", err.Error())
	}

	return &pb.UpdateAppointmentStatusResponse{}, nil
}
