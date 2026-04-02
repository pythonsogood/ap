package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/handler"
)

func SetupDoctorRoutes(engine *gin.Engine, doctor_handler handler.DoctorHandler) {
	doctor_group := engine.Group("/doctors")

	doctor_group.GET("", doctor_handler.GETList)
	doctor_group.GET("/:id", doctor_handler.GETByID)
	doctor_group.POST("", doctor_handler.POST)
}
