.DELETE_ON_ERROR: clean

EXECUTABLES = go
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

BUILD_DIR := ./build
PROJECTS_PATH := $(shell ls -d cmd/*)
PROJECTS_NAME := $(foreach dir_name, $(PROJECTS_PATH), $(shell basename $(dir_name)) )
PROJECT_DEPENDENCIES := $(shell go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all)


all: clean test build

mod-update: tidy
	$(foreach dep, $(PROJECT_DEPENDENCIES), $(shell go get -u $(dep)))
	go mod tidy

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
	$(foreach proj_name, $(PROJECTS_NAME), $(shell CGO_ENABLED=0 go build -o ./$(BUILD_DIR)/$(proj_name) ./cmd/$(proj_name)/))

clean:
	rm -rf $(BUILD_DIR) ./*.out