# Design document

Felek runs arbitrary Unix processes ("jobs") on the server machine, allowing a client to query job status, stop running jobs, and inspect its stdout and stderr streams ("logs").

## Architecture

Felek uses Go's `exec.Cmd` to run external processes, and provides a gRPC interface for clients.

When a `Start` request is received, Felek will try to start an external process using the provided settings (command path and arguments). If the process was started successfully, Felek will generate a unique UUID for it an store it in a key-value store for later queries.

The key-value store for started jobs is protected by an `RWMutex` to allow for many readers or a single writer. `Start` and `Stop` require a write lock, while other RPCs require a read lock.

For each created job, Felek creates two files to store the process's stdout and stderr streams, respectively. This allows easily querying existing logs as well as "following" logs for streaming, for example using existing wrapper around `inotify` such as [`hpcloud/tail`](https://github.com/hpcloud/tail). Felek starts a goroutine for each job to `.Wait()` on the job to complete and then close the file handles to its log files.

For querying logs, Felek use gRPC's streaming functionality to stream logs line-by-line, optionally continuing to stream after the existing content of the log file is exhaust. Logs are assumed to consist of lines of UTF-8 encoded text.

Connections from the client to the server will be encrypted and authenticated using mTLS.

## mTLS setup

### Authentication

We create a self-signed root X.509 CA and an intermediate CA signed by the root.
The intermediate CA is then used to sign both server and client certificates.
The client and server then use the intermediate CA to verify each other.

All communication will use TLS 1.3 and support the cipher suites offered for TLS 1.3 by [`crypto/tls`](https://golang.org/pkg/crypto/tls/) (currently `TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384` and `TLS_CHACHA20_POLY1305_SHA256`, in that order). Other configuration settings will use the `crypto/tls` [defaults](https://golang.org/pkg/crypto/tls/#Config).

### Authorization

We authorize the user to either "user" or "admin" privileges based on the CN section of the client certificate.


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
}

message LogsRequest {
    JobID id = 1; // ID of the job whose logs we want to get
    bool follow = 2; // Keep streaming the logs as they come in
    bool tail = 3; // If `follow=true`, don't stream existing lines before following
}

message JobStatus {
    JobID id = 1; // A unique ID
    oneof jobStatus {
        RunningJob runningJob = 2;
        StoppedJob stoppedJob = 3;
    }
}

message RunningJob {
    // XXX: Don't know if we only want this available for running processes,
    //      but otherwise I'm not sure what I can put in the RunningJob message
    int64 pid = 1; // PID of the job's root process
}

message StoppedJob {
    oneof exit {
        int64 exitCode = 1; // The exit code of the job's process
        string exitError = 2; // The error string in case a job's processes was killed by a signal
    }
    int64 systemTime = 3; // Elapsed system time, in nanoseconds
    int64 userTime = 4; // Elapsed user time, in nanoseconds
    bool stopped = 5; // Was the job stopped by a user
}

// A unique ID of a started job (a UUID string)
message JobID {
    string value = 1;
}

// Represents a single line of logs
message LogLine {
    string value = 1;
}
```
