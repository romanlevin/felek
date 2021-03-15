package main

import (
	pb "github.com/romanlevin/felek/jobs"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const gracefulStopTimeout = 10 * time.Second

func main() {
	// Catch SIGINT and SIGTERM to shut down gracefully
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	lis, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	jobServer, err := newServer()
	if err != nil {
		log.Fatalf("error creating jobs server: %s", err.Error())
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterJobsServer(gRPCServer, jobServer)

	stopped := make(chan struct{})
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		close(stopped)
	}()

	log.Print("Serving gRPC service")

	// Wait for an interrupt
	<-signals

	log.Print("Shutting down server")

	// Start graceful shutdown
	go func() {
		gRPCServer.GracefulStop()
	}()

	// Force shutdown if GracefulStop is taking a while
	t := time.NewTimer(gracefulStopTimeout)
	select {
	case <-t.C:
		log.Print("Timed out, forcing a shutdown")
		gRPCServer.Stop()
	case <-stopped:
		t.Stop()
		log.Print("Server shut down gracefully")
	}
}
