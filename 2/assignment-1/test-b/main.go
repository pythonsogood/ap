package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/pythonsogood/test_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewBundle().TransportCredentials()))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	c := pb.NewAClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World!"})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(r.Message)
}
