package grpc

import (
	"github.com/pythonsogood/ap-assignment1/doctor/internal/handler"
	pb "github.com/pythonsogood/ap-assignment1/proto"
	"google.golang.org/grpc"
)

func SetupDoctorTransport(server *grpc.Server, doctor_handler *handler.DoctorGRPCHandler) error {
	pb.RegisterDoctorServiceServer(server, doctor_handler)

	return nil
}
