# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGENERATE=$(GOCMD) generate
BINARY_NAME=passa
BUILD_DIR=build
BINARY_UNIX=$(BINARY_NAME)_unix

compile:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)

validate:
	ruby -e "require 'yaml';YAML.load_file('./passa-states.yml')"

test:
	$(GOGENERATE)
	$(GOTEST) -v

cover:
	$(GOTEST) -coverprofile cp.out
	$(GOCMD) tool cover -html=cp.out