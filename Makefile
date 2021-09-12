.DELETE_ON_ERROR:

EXECUTABLES = go
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

PROJECTS_PATH := $(shell ls -d cmd/*)
PROJECTS_NAME := $(foreach dir_name, $(PROJECTS_PATH), $(shell basename $(dir_name)) )


all: clean test build

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

generate:
	go generate ./...

test: tidy fmt vet generate
	go test -race -covermode=atomic -coverprofile coverage.out -tags=unit ./...

test-coverage: test
	go tool cover -html=coverage.out

build:
	$(foreach proj_name, $(PROJECTS_NAME), $(shell CGO_ENABLED=0 go build -o ./bin/$(proj_name) ./cmd/$(proj_name)/))

clean:
	rm -rf ./bin ./*.out