package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/pythonsogood/test_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAServer
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("got request")
	return &pb.HelloResponse{Message: fmt.Sprintf("Hello %s", in.GetName())}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")

	if err != nil {
		panic(err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterAServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		panic(err.Error())
	}
}
