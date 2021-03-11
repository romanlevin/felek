package main

import (
	"context"
	"fmt"
	pb "github.com/romanlevin/felek/jobs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
)

type server struct {
	pb.UnimplementedJobsServer
	jobsLock *sync.RWMutex
	jobs     map[string]*job
	logDir   string
}

func newServer() *server {
	logDir, err := ioutil.TempDir("", "felek-logs-*")
	if err != nil {
		panic(err)
	}

	s := &server{jobs: make(map[string]*job), logDir: logDir, jobsLock: &sync.RWMutex{}}
	return s
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

// waitOnJob is meant to be run as a goroutine
// It captures the process's exit error in case it is killed by a signal
// Calling Wait on the exec.Cmd is also necessary to update ProcessState
func (s *server) waitOnJob(j *job) {
	err := j.cmd.Wait()
	if err != nil {
		s.jobsLock.Lock()
		defer s.jobsLock.Unlock()
		j.exitError = err
		log.Printf("job %v exited with error %#v", j.id, j.exitError)
	}
	log.Printf("job %v exited with exit code %v", j.id, j.cmd.ProcessState.ExitCode())

	// TODO close log file handlers
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
	go s.waitOnJob(job)
	return jobStatus(job), nil
}

func (s *server) Stop(ctx context.Context, id *pb.JobID) (*pb.JobStatus, error) {
	jobID := id.Value
	j, ok := s.getJob(jobID)
	if !ok {
		return nil, fmt.Errorf("no such job id %v", jobID)
	}

	err := j.cmd.Process.Kill()
	if err != nil {
		return jobStatus(j), fmt.Errorf("failed to kill: %w", err)
	}

	s.jobsLock.Lock()
	defer s.jobsLock.Unlock()
	j.stopped = true

	return jobStatus(j), nil
}

func (s *server) Status(ctx context.Context, id *pb.JobID) (*pb.JobStatus, error) {
	jobID := id.Value
	j, ok := s.getJob(jobID)
	if !ok {
		return nil, fmt.Errorf("no such job id %v", jobID)
	}

	return jobStatus(j), nil
}

func (s *server) Stdout(request *pb.LogsRequest, logsServer pb.Jobs_StdoutServer) error {
	panic("implement me")
}

func (s *server) Stderr(request *pb.LogsRequest, logsServer pb.Jobs_StderrServer) error {
	panic("implement me")
}
