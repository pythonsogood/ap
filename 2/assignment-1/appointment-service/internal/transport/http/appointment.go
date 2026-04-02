package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/handler"
)

func SetupAppointmentTransport(engine *gin.Engine, appointment_handler handler.AppointmentHandler) error {
	appointment_group := engine.Group("/appointments")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("validstatus", appointment_handler.StatusValidator); err != nil {
			return err
		}
	}

	appointment_group.GET("/", appointment_handler.GETList)
	appointment_group.GET("/:id", appointment_handler.GETByID)
	appointment_group.POST("/", appointment_handler.POST)
	appointment_group.PATCH("/:id/status", appointment_handler.PATCHStatusByID)

	return nil
}
