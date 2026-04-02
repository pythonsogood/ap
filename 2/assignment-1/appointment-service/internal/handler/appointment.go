package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/service"
)

type AppointmentHandler interface {
	GETByID(c *gin.Context)
	GETList(c *gin.Context)
	POST(c *gin.Context)
	PATCHStatusByID(c *gin.Context)

	StatusValidator(fl validator.FieldLevel) bool
}

type appointmentHandlerImpl struct {
	service service.AppointmentService
}

func NewAppointmentHandler(service service.AppointmentService) AppointmentHandler {
	return &appointmentHandlerImpl{
		service: service,
	}
}

func (h *appointmentHandlerImpl) GETByID(c *gin.Context) {
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

func (h *appointmentHandlerImpl) GETList(c *gin.Context) {
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

func (h *appointmentHandlerImpl) POST(c *gin.Context) {
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

	c.JSON(http.StatusOK, appointment)
}

type AppointmentPATCHStatusBind struct {
	Status model.Status `json:"status" binding:"required,validstatus"`
}

func (h *appointmentHandlerImpl) PATCHStatusByID(c *gin.Context) {
	id := c.Param("id")

	var status_bind AppointmentPATCHStatusBind

	if err := c.ShouldBindJSON(&status_bind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(id, status_bind.Status); err != nil {
		status_code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			status_code = http.StatusNotFound
		}

		c.JSON(status_code, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *appointmentHandlerImpl) StatusValidator(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(model.Status)

	if !ok {
		return false
	}

	return status == model.StatusDone || status == model.StatusInProgress || status == model.StatusNew
}
