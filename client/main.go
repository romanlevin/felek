package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/romanlevin/felek/jobs"
	"google.golang.org/grpc"
)

const (
	address = "localhost:12345"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := jobs.NewJobsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch command := os.Args[1]; command {
	case "run":
		runJob(c, ctx)
	case "status":
		statusJob(c, ctx)
	case "stop":
		stopJob(c, ctx)
	default:
		log.Fatalf("must run with %s [run|status|stop] ..., instead got %s", os.Args[0], os.Args[1])
	}
}

func runJob(c jobs.JobsClient, ctx context.Context) {
	request := &jobs.JobStartRequest{
		Path: os.Args[2],
		Args: os.Args[3:],
	}
	log.Printf("sending request %v", request)
	r, err := c.Start(ctx, request)
	if err != nil {
		log.Fatalf("could not start: %#v", err.Error())
	}
	log.Printf("Job name: %v", r)
}

func stopJob(c jobs.JobsClient, ctx context.Context) {
	id := &jobs.JobID{Value: os.Args[2]}
	status, err := c.Stop(ctx, id)
	if err != nil {
		log.Fatalf("could not kill: %#v", err.Error())
	}
	logJobStatus(status)
}

func statusJob(c jobs.JobsClient, ctx context.Context) {
	id := &jobs.JobID{Value: os.Args[2]}
	status, err := c.Status(ctx, id)
	if err != nil {
		log.Fatalf("could not get status: %#v", err.Error())
	}
	logJobStatus(status)
}

func logJobStatus(status *jobs.JobStatus) {
	switch jobState := status.JobState.(type) {
	case *jobs.JobStatus_StoppedJob:
		log.Printf("stopped job(%s): %v, stopped by user: %t", status.Id.Value, jobState.StoppedJob, jobState.StoppedJob.Stopped)
	case *jobs.JobStatus_RunningJob:
		log.Printf("running job(%s): %v", status.Id.Value, jobState.RunningJob)
	default:
		log.Fatalf("JobStatus.JobState has unexpected type %T", jobState)

	}
}
