package main

import (
	pb "github.com/romanlevin/felek/jobs"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":12345")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRPCServer := grpc.NewServer()

	jobServer, err := newServer()
	if err != nil {
		log.Fatalf("error creating jobs server: %s", err.Error())
	}

	pb.RegisterJobsServer(gRPCServer, jobServer)
	log.Print("Serving gRPC service")
	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
