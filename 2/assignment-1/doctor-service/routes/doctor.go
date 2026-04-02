package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/handlers"
)

func SetupDoctorRoutes(engine *gin.Engine) {
	doctor_group := engine.Group("/doctors")

	doctor_group.GET("", handlers.DoctorGETListHandler)
	doctor_group.GET("/:id", handlers.DoctorGETByIDHandler)
	doctor_group.POST("", handlers.DoctorPOSTHandler)
}
