package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/appointment/cmd/appointment/config"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/database"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/handler"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/service"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/transport/http"
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

	appointment_db, err := database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

	if err != nil {
		panic(err.Error())
	}

	if err := database.SQLiteInitDB(appointment_db, []database.Model{&model.Appointment{}}); err != nil {
		panic(err.Error())
	}

	appointment_repo := repository.NewSQLiteAppointmentRepository(appointment_db)
	appointment_service := service.NewAppointmentService(appointment_repo)
	appointment_handler := handler.NewAppointmentHandler(appointment_service)

	if err := http.SetupAppointmentTransport(router, appointment_handler); err != nil {
		panic(err.Error())
	}

	if err := router.Run(server_addr); err != nil {
		panic(err.Error())
	}
}
