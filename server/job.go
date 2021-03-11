package main

import (
	"github.com/google/uuid"
	pb "github.com/romanlevin/felek/jobs"
	"log"
	"os/exec"
	"sync"
)

type job struct {
	id        string    // A UUID string
	cmd       *exec.Cmd // The underlying Cmd object wrapping the job process
	stopped   bool      // Has this job been stopped by a client?
	exitError error     // The error returned by Wait on the Cmd
	owner     string    // The username of the user who started the job
	lock      *sync.RWMutex
}

// exitedWithExitCode assumes a lock is held
func (j *job) exitedWithExitCode() bool {
	return j.exitError == nil && j.cmd.ProcessState != nil
}

// running assumes a lock is held
func (j *job) running() bool {
	return j.exitError == nil && j.cmd.ProcessState == nil
}

func (j *job) stop() error {
	j.lock.Lock()
	defer j.lock.Unlock()

	if err := j.cmd.Process.Kill(); err != nil{
		return err
	}

	// XXX: Kind of a philosophical question - if the user tries to stop the process but that fails,
	// do we count the job as stopped by the user? For now, no.
	j.stopped = true

	return nil
}

func (j *job) status() *pb.JobStatus {
	j.lock.RLock()
	defer j.lock.RUnlock()

	cmd := j.cmd
	status := &pb.JobStatus{
		Id: &pb.JobID{Value: j.id},
	}

	if j.exitedWithExitCode() {
		status.JobState = &pb.JobStatus_StoppedJob{
			StoppedJob: &pb.StoppedJob{
				Exit: &pb.StoppedJob_ExitCode{
					ExitCode: int64(cmd.ProcessState.ExitCode()),
				},
				SystemTime: int64(cmd.ProcessState.SystemTime()),
				UserTime:   int64(cmd.ProcessState.UserTime()),
				Stopped:    j.stopped,
			},
		}
		return status
	}

	if j.running() {
		status.JobState = &pb.JobStatus_RunningJob{
			RunningJob: &pb.RunningJob{Pid: int64(j.cmd.Process.Pid)},
		}
		return status
	}

	// If job is not running and has exitCode, it was killed with a signal
	status.JobState = &pb.JobStatus_StoppedJob{
		StoppedJob: &pb.StoppedJob{
			Exit: &pb.StoppedJob_ExitError{
				ExitError: j.exitError.Error(),
			},
			SystemTime: int64(cmd.ProcessState.SystemTime()),
			UserTime:   int64(cmd.ProcessState.UserTime()),
			Stopped:    j.stopped,
		},
	}
	return status
}

// newJob creates a job struct for a new Cmd, assigning it a UUID
func newJob(cmd *exec.Cmd, owner string) *job {
	id := uuid.NewString()
	return &job{
		id:    id,
		cmd:   cmd,
		owner: owner,
		lock: &sync.RWMutex{},
	}
}

// wait is meant to be run as a goroutine
// It captures the process's exit error in case it is killed by a signal
// Calling Wait on the exec.Cmd is also necessary to update ProcessState
func (j *job) wait() {
	err := j.cmd.Wait()
	if err != nil {
		j.lock.Lock()
		defer j.lock.Unlock()
		j.exitError = err
		log.Printf("job %v exited with error %#v", j.id, j.exitError)
	}
	log.Printf("job %v exited with exit code %v", j.id, j.cmd.ProcessState.ExitCode())

	// TODO close log file handlers
}
