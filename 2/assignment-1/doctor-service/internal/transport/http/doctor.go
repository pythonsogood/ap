package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
)

func SetupDoctorTransport(engine *gin.Engine, doctor_handler handler.DoctorHandler) error {
	doctor_group := engine.Group("/doctors")

	doctor_group.GET("/", doctor_handler.GETList)
	doctor_group.GET("/:id", doctor_handler.GETByID)
	doctor_group.POST("/", doctor_handler.POST)

	return nil
}
