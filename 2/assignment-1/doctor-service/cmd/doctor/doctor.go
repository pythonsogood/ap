package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/cmd/doctor/config"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/database"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/model"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/repository"
	"github.com/pythonsogood/ap-assignment1/doctor/internal/service"
	grpc_transport "github.com/pythonsogood/ap-assignment1/doctor/internal/transport/grpc"
	http_transport "github.com/pythonsogood/ap-assignment1/doctor/internal/transport/http"
	"google.golang.org/grpc"
)

func serve_http(server_addr string, doctor_service service.DoctorService) (*gin.Engine, handler.DoctorHTTPHandler, error) {
	router := gin.Default()

	doctor_handler := handler.NewDoctorHTTPHandler(doctor_service)

	if err := http_transport.SetupDoctorTransport(router, doctor_handler); err != nil {
		return router, doctor_handler, err
	}

	if err := router.Run(server_addr); err != nil {
		return router, doctor_handler, err
	}

	return router, doctor_handler, nil
}

func serve_grpc(server_addr string, doctor_service service.DoctorService) (*net.Listener, *handler.DoctorGRPCHandler, error) {
	lis, err := net.Listen("tcp", server_addr)

	if err != nil {
		return &lis, nil, err
	}

	doctor_handler := handler.NewDoctorGRPCHandler(doctor_service)

	s := grpc.NewServer()

	err = grpc_transport.SetupDoctorTransport(s, doctor_handler)

	if err != nil {
		return &lis, doctor_handler, err
	}

	if err := s.Serve(lis); err != nil {
		return &lis, doctor_handler, err
	}

	return &lis, doctor_handler, nil
}

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

	doctor_db, err := database.SQLiteConnectDB(conf.Database.Sqlite3.Source)

	if err != nil {
		panic(err.Error())
	}

	if err := database.SQLiteInitDB(doctor_db, []database.Model{&model.Doctor{}}); err != nil {
		panic(err.Error())
	}

	doctor_repo := repository.NewSQLiteDoctorRepository(doctor_db)
	doctor_service := service.NewDoctorService(doctor_repo)

	// _, _, err = serve_http(server_addr, doctor_service)

	// if err != nil {
	// 	panic(err.Error())
	// }

	_, _, err = serve_grpc(server_addr, doctor_service)

	if err != nil {
		panic(err.Error())
	}
}
