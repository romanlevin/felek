.PHONY: clean server

all: clean protoc server

protoc: jobs/felek.proto
	protoc \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=. \
	  --go-grpc_opt=paths=source_relative \
	  jobs/felek.proto

server:
	go run ./server

clean:
	rm jobs/*.go