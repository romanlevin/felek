# felek

## Usage

To generate protobuf and gRPC code for Go, run

```shell
$ make protoc
```

To run the server, run

```shell
$ make server
```

Examples for using the CLI client:

```shell
# Run job
$ go run ./client run bash -c "sleep 30 && echo something"
# Mixed server and client output
2021/03/10 18:45:44 sending request path:"bash" args:"-c" args:"sleep 30 && echo something"
2021/03/10 18:45:45 Started job (be70da72-6abb-498a-b456-19daba3e2b82): "/usr/bin/bash" []string{"bash", "-c", "sleep 30 && echo something"}
2021/03/10 18:45:45 Job name: id:{value:"be70da72-6abb-498a-b456-19daba3e2b82"} runningJob:{pid:31763}

# Check status of running job
$ go run ./client status be70da72-6abb-498a-b456-19daba3e2b82
2021/03/10 18:45:53 running job(be70da72-6abb-498a-b456-19daba3e2b82): pid:31763

# Job finishes
something
2021/03/10 18:46:15 job be70da72-6abb-498a-b456-19daba3e2b82 exited with exit code 0

# Check status of finished job
$ go run ./client status be70da72-6abb-498a-b456-19daba3e2b82
2021/03/10 18:46:20 stopped job(be70da72-6abb-498a-b456-19daba3e2b82): exitCode:0 userTime:2099000, stopped by user: false

# Check on non-existent jbo
$ go run ./client status blah
2021/03/10 18:50:14 could not get status: "rpc error: code = Unknown desc = no such job id blah"
exit status 1

# Start and stop a job
$ go run ./client run bash -c "sleep 600 && echo something"
# Mixed server and client output
2021/03/10 18:44:38 sending request path:"bash" args:"-c" args:"sleep 600 && echo something"
2021/03/10 18:44:38 Started job (6648ae22-5c7f-4ea6-b263-f38e25b4de68): "/usr/bin/bash" []string{"bash", "-c", "sleep 600 && echo something"}
2021/03/10 18:44:38 Job name: id:{value:"6648ae22-5c7f-4ea6-b263-f38e25b4de68"} runningJob:{pid:30976}

# Stop the job
$ go run ./client stop 6648ae22-5c7f-4ea6-b263-f38e25b4de68
# Server output
2021/03/10 18:44:53 job 6648ae22-5c7f-4ea6-b263-f38e25b4de68 exited with error "signal: killed"
2021/03/10 18:44:53 job 6648ae22-5c7f-4ea6-b263-f38e25b4de68 exited with exit code -1
# Client output (job has not updated yet, displayed as running)
# XXX: Prevent this?
2021/03/10 18:44:53 running job(6648ae22-5c7f-4ea6-b263-f38e25b4de68): pid:30976

# Check status of stopped job
$ go run ./client status 6648ae22-5c7f-4ea6-b263-f38e25b4de68
2021/03/10 18:45:16 stopped job(6648ae22-5c7f-4ea6-b263-f38e25b4de68): exitError:"signal: killed" userTime:1097000 stopped:true, stopped by user: true
```