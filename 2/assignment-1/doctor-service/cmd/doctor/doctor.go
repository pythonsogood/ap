package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/cmd/doctor/config"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/database"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/repository"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/route"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
)

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	switch conf.Database.Type {
	case config.DatabaseTypeSQLite:
	default:
		panic("Unsupported database type!")
	}

	server_addr := fmt.Sprintf(":%d", conf.Server.Port)

	router := gin.Default()

	doctor_db, err := database.NewSQLiteDB(conf.Database.Sqlite3.Source)

	if err != nil {
		panic(err.Error())
	}

	doctor_repo := repository.NewSQLiteDoctorRepository(doctor_db)
	doctor_service := service.NewDoctorService(doctor_repo)
	doctor_handler := handler.NewDoctorHandler(doctor_service)

	route.SetupDoctorRoutes(router, doctor_handler)

	if err := router.Run(server_addr); err != nil {
		panic(err.Error())
	}
}
