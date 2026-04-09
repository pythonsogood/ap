package grpc

import (
	"github.com/pythonsogood/ap-assignment1/appointment/internal/handler"
	pb "github.com/pythonsogood/ap-assignment1/proto"
	"google.golang.org/grpc"
)

func SetupAppointmentdTransport(server *grpc.Server, appointment_handler *handler.AppointmentGRPCHandler) error {
	pb.RegisterAppointmentServiceServer(server, appointment_handler)

	return nil
}
