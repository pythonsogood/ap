package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/database"
	"github.com/pythonsogood/ap-assignment1/doctor/handler"
	"github.com/pythonsogood/ap-assignment1/doctor/repository"
	"github.com/pythonsogood/ap-assignment1/doctor/route"
	"github.com/pythonsogood/ap-assignment1/doctor/service"
)

func main() {
	router := gin.Default()

	doctor_db, err := database.SQLiteConnectDB("doctor-service.db")

	if err != nil {
		panic(err.Error())
	}

	doctor_repo := repository.NewSQLiteDoctorRepository(doctor_db)
	doctor_service := service.NewDoctorService(doctor_repo)
	doctor_handler := handler.NewDoctorHandler(&doctor_service)

	route.SetupDoctorRoutes(router, doctor_handler)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})

	if err := router.Run(":8081"); err != nil {
		panic(err.Error())
	}
}
