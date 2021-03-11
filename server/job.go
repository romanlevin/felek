package main

import (
	"github.com/google/uuid"
	pb "github.com/romanlevin/felek/jobs"
	"os/exec"
)

type job struct {
	id        string    // A UUID string
	cmd       *exec.Cmd // The underlying Cmd object wrapping the job process
	stopped   bool      // Has this job been stopped by a client?
	exitError error     // The error returned by Wait on the Cmd
	owner     string    // The username of the user who started the job
}

func (j *job) exitedWithExitCode() bool {
	return j.exitError == nil && j.cmd.ProcessState != nil
}

func (j *job) running() bool {
	return j.exitError == nil && j.cmd.ProcessState == nil
}

func (j *job) killedWithSignal() bool {
	return j.exitError != nil
}

func (j *job) status() *pb.JobStatus {
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
	}
}
