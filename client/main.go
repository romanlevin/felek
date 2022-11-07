package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/romanlevin/felek/jobs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:12345"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := jobs.NewJobsClient(conn)

	ctx := context.Background()

	if len(os.Args) < 3 {
		log.Fatalf("expecting command argument")
	}

	switch command := os.Args[1]; command {
	case "run":
		err = runJob(ctx, client)
	case "status":
		err = statusJob(ctx, client)
	case "stop":
		err = stopJob(ctx, client)
	default:
		log.Fatalf("must run with %s [run|status|stop] ..., instead got %s", os.Args[0], os.Args[1])
	}

	if err != nil {
		log.Fatalf("error encountered: %s", err.Error())
	}
}

func runJob(ctx context.Context, client jobs.JobsClient) error {
	request := &jobs.JobStartRequest{
		Path: os.Args[2],
		Args: os.Args[3:],
	}
	log.Printf("sending request %v", request)
	status, err := client.Start(ctx, request)
	if err != nil {
		return fmt.Errorf("could not start: %w", err)
	}
	if err = logJobStatus(status); err != nil {
		return err
	}
	return nil
}

func stopJob(ctx context.Context, client jobs.JobsClient) error {
	id := &jobs.JobID{Value: os.Args[2]}
	status, err := client.Stop(ctx, id)
	if err != nil {
		return fmt.Errorf("could not kill: %w", err)
	}
	if err = logJobStatus(status); err != nil {
		return err
	}
	return nil
}

func statusJob(ctx context.Context, client jobs.JobsClient) error {
	id := &jobs.JobID{Value: os.Args[2]}
	status, err := client.Status(ctx, id)
	if err != nil {
		return fmt.Errorf("could not get status: %w", err)
	}
	if err = logJobStatus(status); err != nil {
		return err
	}
	return nil
}

func logJobStatus(status *jobs.JobStatus) error {
	switch jobState := status.JobState.(type) {
	case *jobs.JobStatus_StoppedJob:
		log.Printf("stopped job(%s): %v, stopped by user: %t", status.Id.Value, jobState.StoppedJob, jobState.StoppedJob.Stopped)
	case *jobs.JobStatus_RunningJob:
		log.Printf("running job(%s): %v", status.Id.Value, jobState.RunningJob)
	default:
		return fmt.Errorf("JobStatus.JobState has unexpected type %T", jobState)
	}
	return nil
}
