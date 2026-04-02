package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/cmd/doctor/config"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/database"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/repository"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/transport/http"
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

	doctor_db, err := database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

	if err != nil {
		panic(err.Error())
	}

	if err := database.SQLiteInitDB(doctor_db, []database.Model{&model.Doctor{}}); err != nil {
		panic(err.Error())
	}

	doctor_repo := repository.NewSQLiteDoctorRepository(doctor_db)
	doctor_service := service.NewDoctorService(doctor_repo)
	doctor_handler := handler.NewDoctorHandler(doctor_service)

	if err := http.SetupDoctorTransport(router, doctor_handler); err != nil {
		panic(err.Error())
	}

	if err := router.Run(server_addr); err != nil {
		panic(err.Error())
	}
}
