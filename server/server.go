package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	pb "github.com/romanlevin/felek/jobs"
	"golang.org/x/sync/errgroup"
)

type server struct {
	pb.UnimplementedJobsServer
	jobs     map[string]*job
	jobsLock sync.RWMutex
	logDir   string
}

func newServer() (*server, error) {
	logDir, err := os.CreateTemp("", "felek-logs-*")
	if err != nil {
		return nil, err
	}

	s := &server{jobs: make(map[string]*job), logDir: logDir.Name(), jobsLock: sync.RWMutex{}}
	return s, nil
}

func (s *server) storeJob(j *job) {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()
	s.jobs[j.id] = j
}

func (s *server) getJob(id string) (*job, bool) {
	s.jobsLock.RLock()
	defer s.jobsLock.RUnlock()
	j, ok := s.jobs[id]
	return j, ok

}

func (s *server) Start(ctx context.Context, request *pb.JobStartRequest) (*pb.JobStatus, error) {
	cmd := exec.Command(request.Path, request.Args...)
	// TODO: Pipe stdout and stderr to files
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	// TODO: Assign owner using the CN field of the client certificate
	job := newJob(cmd, "")
	log.Printf("Started job (%v): %#v %#v", job.id, job.cmd.Path, job.cmd.Args)
	s.storeJob(job)
	go job.wait()
	return job.status(), nil
}

func (s *server) Stop(ctx context.Context, id *pb.JobID) (*pb.JobStatus, error) {
	jobID := id.Value
	job, ok := s.getJob(jobID)
	if !ok {
		return nil, fmt.Errorf("no such job id %v", jobID)
	}

	err := job.stop()
	if err != nil {
		return nil, fmt.Errorf("failed to stop: %w", err)
	}

	return job.status(), nil
}

func (s *server) ShutDown() error {
	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()
	g := errgroup.Group{}
	for _, job := range s.jobs {
		g.Go(job.stop)
	}
	return g.Wait()
}

func (s *server) Status(ctx context.Context, id *pb.JobID) (*pb.JobStatus, error) {
	jobID := id.Value
	job, ok := s.getJob(jobID)
	if !ok {
		return nil, fmt.Errorf("no such job id %v", jobID)
	}

	return job.status(), nil
}

func (s *server) Stdout(request *pb.LogsRequest, logsServer pb.Jobs_StdoutServer) error {
	panic("implement me")
}

func (s *server) Stderr(request *pb.LogsRequest, logsServer pb.Jobs_StderrServer) error {
	panic("implement me")
}
