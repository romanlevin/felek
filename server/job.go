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

func jobStatus(j *job) *pb.JobStatus {
	c := j.cmd
	if j.exitError == nil && c.ProcessState != nil {
		// Job exited without exitError
		return &pb.JobStatus{
			Id: &pb.JobID{Value: j.id},
			JobState: &pb.JobStatus_StoppedJob{
				StoppedJob: &pb.StoppedJob{
					Exit: &pb.StoppedJob_ExitCode{
						ExitCode: int64(c.ProcessState.ExitCode()),
					},
					SystemTime: int64(c.ProcessState.SystemTime()),
					UserTime:   int64(c.ProcessState.UserTime()),
					Stopped:    j.stopped,
				},
			},
		}
	} else if c.ProcessState == nil {
		// Job still running
		return &pb.JobStatus{
			Id: &pb.JobID{Value: j.id},
			JobState: &pb.JobStatus_RunningJob{
				RunningJob: &pb.RunningJob{Pid: int64(j.cmd.Process.Pid)},
			},
		}
	} else {
		// Job exited with exitError
		return &pb.JobStatus{
			Id: &pb.JobID{Value: j.id},
			JobState: &pb.JobStatus_StoppedJob{
				StoppedJob: &pb.StoppedJob{
					Exit: &pb.StoppedJob_ExitError{
						ExitError: j.exitError.Error(),
					},
					SystemTime: int64(c.ProcessState.SystemTime()),
					UserTime:   int64(c.ProcessState.UserTime()),
					Stopped:    j.stopped,
				},
			},
		}
	}
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
