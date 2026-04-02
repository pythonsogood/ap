package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
)

type DoctorHandler interface {
	GETByID(c *gin.Context)
	GETList(c *gin.Context)
	POST(c *gin.Context)
}

type doctorHandlerImpl struct {
	service service.DoctorService
}

func NewDoctorHandler(service *service.DoctorService) DoctorHandler {
	return &doctorHandlerImpl{}
}

func (h *doctorHandlerImpl) GETByID(c *gin.Context) {
	id := c.Param("id")

	doctor, err := h.service.GetDoctor(id)

	if err != nil {
		c.Error(err)

		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *doctorHandlerImpl) GETList(c *gin.Context) {
	doctors, err := h.service.GetAllDoctors()

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, doctors)
}

type DoctorPOSTBind struct {
	FullName       string `json:"full_name" binding:"required"`
	Specialization string `json:"specialization" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

func (h *doctorHandlerImpl) POST(c *gin.Context) {
	var doctor_bind DoctorPOSTBind

	if err := c.BindJSON(&doctor_bind); err != nil {
		return
	}

	doctor, err := h.service.CreateDoctor(doctor_bind.FullName, doctor_bind.Specialization, doctor_bind.Email)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, doctor)
}
