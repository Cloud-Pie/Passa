# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGENERATE=$(GOCMD) generate
BINARY_NAME=passa
BUILD_DIR=build
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all test clean server

compile:
	$(GOBUILD)

validate:
	$(GOGENERATE) #this also validates the pass-states.yml

test:
	go generate
	go clean -testcache
	go test ./...  -cover

cover:
	$(GOTEST) -coverprofile cp.out
	$(GOCMD) tool cover -html=cp.out

run:
	$(GOGENERATE)
	$(GOBUILD)
	./Passa

clean:
	$(GOCLEAN) -testcache
	- rm .db.json
	- rm  build/*
	

server:
	make validate
	make compile
	./Passa --no-cloud

dist:
	env GOOS=linux GOARCH=amd64 go build -o build/passa_linux
	go build -o build/passa_mac

linux:
	env GOOS=linux GOARCH=amd64 go build -o build/passa_linux
	scp build/passa_linux  root@10.155.209.58:./Passa/.
