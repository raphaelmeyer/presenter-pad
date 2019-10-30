all: build

build:
	go build ./...

install:
	go install ./...

run:
	go run cmd/presenter-pad/main.go

install-tools:
	go mod download
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
