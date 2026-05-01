package main

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/appointment/cmd/appointment/config"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/database"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/event"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/handler"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/model"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/repository"
	"github.com/pythonsogood/ap-assignment1/appointment/internal/service"
	grpc_transport "github.com/pythonsogood/ap-assignment1/appointment/internal/transport/grpc"
	http_transport "github.com/pythonsogood/ap-assignment1/appointment/internal/transport/http"
	"google.golang.org/grpc"
)

func serve_http(server_addr string, appointment_service service.AppointmentService, event_publisher event.EventPublisher) (*gin.Engine, handler.AppointmentHTTPHandler, error) {
	router := gin.Default()

	appointment_handler := handler.NewAppointmentHTTPHandler(appointment_service, event_publisher)

	if err := http_transport.SetupAppointmentTransport(router, appointment_handler); err != nil {
		return router, appointment_handler, err
	}

	if err := router.Run(server_addr); err != nil {
		return router, appointment_handler, err
	}

	return router, appointment_handler, nil
}

func serve_grpc(server_addr string, appointment_service service.AppointmentService, event_publisher event.EventPublisher) (*net.Listener, *handler.AppointmentGRPCHandler, error) {
	lis, err := net.Listen("tcp", server_addr)

	if err != nil {
		return &lis, nil, err
	}

	appointment_handler := handler.NewAppointmentGRPCHandler(appointment_service, event_publisher)

	s := grpc.NewServer()

	err = grpc_transport.SetupAppointmentdTransport(s, appointment_handler)

	if err != nil {
		return &lis, appointment_handler, err
	}

	if err := s.Serve(lis); err != nil {
		return &lis, appointment_handler, err
	}

	return &lis, appointment_handler, nil
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	var appointment_db *sql.DB
	var appointment_repo repository.AppointmentRepository
	var event_publisher event.EventPublisher

	switch conf.Database.Type {
	case config.DatabaseTypeSQLite:
		appointment_db, err = database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

		if err != nil {
			panic(err.Error())
		}

		if err := database.SqlInitModels(appointment_db, []database.Model{&model.Appointment{}}); err != nil {
			panic(err.Error())
		}

		appointment_repo = repository.NewSQLiteAppointmentRepository(appointment_db)
	case config.DatabaseTypePostgres:
		appointment_db, err = database.PostgresConnectDB(conf.Database.Postgres.ConnectionUrl)

		if err != nil {
			panic(err.Error())
		}

		if err := database.SqlInitModels(appointment_db, []database.Model{&model.Appointment{}}); err != nil {
			panic(err.Error())
		}

		appointment_repo = repository.NewPostgresAppointmentRepository(appointment_db)
	default:
		panic("Unsupported database type!")
	}

	switch conf.MessageBroker.Type {
	case config.MessageBrokerTypeNATS:
		nc, err := nats.Connect(conf.MessageBroker.Nats.ConnectionUrl)

		if err != nil {
			panic(err.Error())
		}

		event_publisher = event.NewNATSEventPublisher(nc)
	default:
		panic("Unsupported message broker type!")
	}

	server_addr := fmt.Sprintf(":%d", conf.Server.Port)

	// http_client := http.Client{
	// 	Timeout: time.Duration(conf.Services.Doctor.Timeout*time.Second),
	// }

	// doctor_service := service.NewHTTPDoctorService(conf.Services.Doctor.Address, &http_client)
	doctor_service := service.NewGRPCDoctorService(conf.Services.Doctor.Address, time.Duration(conf.Services.Doctor.Timeout)*time.Second)

	appointment_service := service.NewAppointmentService(appointment_repo, doctor_service)

	// _, _, err = serve_http(server_addr, appointment_service, event_publisher)

	// if err != nil {
	// 	panic(err.Error())
	// }

	_, _, err = serve_grpc(server_addr, appointment_service, event_publisher)

	if err != nil {
		panic(err.Error())
	}
}
