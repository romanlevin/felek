# Design document

Felek runs arbitrary Unix processes ("jobs") on the server machine, allowing a client to query job status, stop running jobs, and inspect its stdout and stderr streams ("logs").

## Architecture

Felek uses Go's `exec.Cmd` to run external processes, and provides a gRPC interface for clients.

When a `Start` request is received, Felek will try to start an external process using the provided settings (command path, arguments, environment variables, and a working directory path). If the process was started successfully, Felek will generate a unique ID for it an store it in a key-value store for later queries.

The key-value store for started jobs is protected by an `RWMutex` to allow for many readers or a single writer. `Start` and `Stop` require a write lock, while other RPCs require a read lock.

For each created job, Felek creates two files to store the process's stdout and stderr streams, respectively. This allows easily querying existing logs as well as "following" logs for streaming, for example using existing wrapper around `inotify` such as [`hpcloud/tail`](https://github.com/hpcloud/tail). Felek starts a goroutine for each job to `.Wait()` on the job to complete and then close the file handles to its log files.

For querying logs, Felek use gRPC's streaming functionality to stream logs line-by-line, optionally continuing to stream after the existing content of the log file is exhaust. Logs are assumed to consist of lines of UTF-8 encoded text.

Connections from the client to the server will be encrypted and authenticated using mTLS. 


## gRPC spec

```protobuf
syntax = "proto3";

service Jobs {
    rpc Start (JobStartRequest) returns (JobStatus);
    rpc Stop (JobID) returns (JobStatus);
    rpc Status (JobID) returns (JobStatus);
    rpc Stdout (LogsRequest) returns (stream LogLine);
    rpc Stderr(LogsRequest) returns (stream LogLine);
}

message JobStartRequest {
    string path = 1;  // Path to executable
    repeated string args = 2; // Represents argv, passed to the executable
    repeated string env = 3; // A list of `key=value` representations of environment variables to set for the process
    string dir = 4; // The working directory
}

message LogsRequest {
    JobID id = 1; // ID of the job whose logs we want to get
    bool follow = 2; // Keep streaming the logs as they come in
    bool tail = 3; // If `follow=true`, don't stream existing lines before following
}

message JobStatus {
    JobID id = 1; // A unique ID
    int64 pid = 2; // PID of the job's process
    bool exited = 3; // Whether the job's process has exited
    // The following are set after a process has exited
    int64 exitCode = 4; // The exit code of the job's process
    int64 systemTime = 5; // Elapsed system time, in nanoseconds
    int64 userTime = 6; // Elapsed user time, in nanoseconds
}

// A unique ID of a started job
message JobID {
    int64 value = 1;
}

// Represents a single line of logs
message LogLine {
    string value = 1;
}
```