.DELETE_ON_ERROR: clean

EXECUTABLES = go
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

PROJECT_NAME      ?= idp-scim-sync
PROJECT_NAMESPACE ?= slashdevops
PROJECT_MODULES_PATH := $(shell ls -d cmd/*)
PROJECT_MODULES_NAME := $(foreach dir_name, $(PROJECT_MODULES_PATH), $(shell basename $(dir_name)) )
PROJECT_DEPENDENCIES := $(shell go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all)

GIT_VERSION  ?= $(shell git rev-parse --abbrev-ref HEAD | cut -d "/" -f 2)
GIT_REVISION ?= $(shell git rev-parse HEAD | tr -d '\040\011\012\015\n')
GIT_BRANCH   ?= $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')
GIT_USER     ?= $(shell git config --get user.name | tr -d '\040\011\012\015\n')
BUILD_DATE   ?= $(shell date +'%Y-%m-%dT%H:%M:%S')

BUILD_DIR       := ./build
DIST_DIR        := ./dist
DIST_ASSEST_DIR := $(DIST_DIR)/assets

GO_LDFLAGS     ?= -ldflags "-X github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Version=$(GIT_VERSION) -X github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Revision=$(GIT_REVISION) -X github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Branch=$(GIT_BRANCH) -X github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.BuildUser=\"$(GIT_USER)\" -X github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.BuildDate=$(BUILD_DATE)"
GO_CGO_ENABLED ?= 0
GO_OPTS        ?= -v
GO_OS          ?= darwin linux windows
GO_ARCH        ?= arm64 amd64
# avoid mocks in tests
GO_FILES       := $(shell go list ./... | grep -v /mocks/)

CONTAINER_OS   ?= linux
CONTAINER_ARCH ?= arm64v8 amd64
CONTAINER_NAMESPACE ?= $(PROJECT_NAMESPACE)
CONTAINER_IMAGE_NAME ?= $(PROJECT_NAME)

DOCKER_CONTAINER_REPO  ?= "docker.io"
AWS_ECR_CONTAINER_REPO ?= "public.ecr.aws/l2n7y5s7"
GITHUB_CONTAINER_REPO  ?= "ghcr.io"

all: clean test build

mod-update: tidy
	$(foreach dep, $(PROJECT_DEPENDENCIES), $(shell go get -u $(dep)))
	go mod tidy

tidy:
	go mod tidy

fmt:
	@go fmt $(GO_FILES)

vet:
	go vet $(GO_FILES)

lint:
	golangci-lint run

generate:
	go generate $(GO_FILES)

test: tidy fmt vet
	go test -race -covermode=atomic -coverprofile coverage.out -tags=unit $(GO_FILES)

test-coverage: test
	go tool cover -html=coverage.out

build:
	$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
		$(shell CGO_ENABLED=$(GO_CGO_ENABLED) go build $(GO_LDFLAGS) $(GO_OPTS) -o ./$(BUILD_DIR)/$(proj_mod) ./cmd/$(proj_mod)/ ) \
	)

build-dist: build
	$(foreach GOOS, $(GO_OS), \
		$(foreach GOARCH, $(GO_ARCH), \
			$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
				$(shell GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(GO_CGO_ENABLED) go build $(GO_LDFLAGS) $(GO_OPTS) -o ./$(DIST_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH)/$(proj_mod) ./cmd/$(proj_mod)/ ) \
			) \
		) \
	)

build-dist-zip:
	mkdir ./$(DIST_ASSEST_DIR);
	$(foreach GOOS, $(GO_OS), \
		$(foreach GOARCH, $(GO_ARCH), \
			$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
				zip --junk-paths -r ./$(DIST_ASSEST_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH).zip ./$(DIST_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH); \
				shasum -a 256 ./$(DIST_ASSEST_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH).zip | cut -d " " -f 1 > ./$(DIST_ASSEST_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH).sha256; \
			) \
		) \
	)

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR) ./*.out .aws-sam/ build.toml

container-build: build-dist
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring $(ARCH), arm64v8), $(eval BIN_ARCH = arm64), $(eval BIN_ARCH = $(ARCH)) ) \
			  echo "Building $(PROJECT_NAME) for OS=$(OS) ARCH=$(ARCH) and BIN_ARCH=$(BIN_ARCH)"; \
				docker build \
					--build-arg ARCH=$(ARCH) \
					--build-arg BIN_ARCH=$(BIN_ARCH) \
					--build-arg OS=$(OS) \
					--platform=$(OS)/$(BIN_ARCH) \
					-t $(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest \
					-t $(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION) \
					-t $(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest \
					-t $(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION) \
					-t $(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest \
					-t $(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION) \
					-t $(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest \
					-t $(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION) \
					./.; \
			))

container-publish-docker: container-build
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring $(ARCH), arm64v8), $(eval BIN_ARCH = arm64),$(eval BIN_ARCH = $(ARCH)) ) \
			docker push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest";  \
			docker push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION)"; \
			))

container-publish-github: container-build
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring $(ARCH), arm64v8), $(eval BIN_ARCH = arm64),$(eval BIN_ARCH = $(ARCH))) \
			docker push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest"; \
			docker push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION)"; \
			))

container-publish-aws-ecr: container-build
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring $(ARCH), arm64v8), $(eval BIN_ARCH = arm64),$(eval BIN_ARCH = $(ARCH))) \
			docker push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):latest"; \
			docker push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME)-$(OS)-$(ARCH):$(GIT_VERSION)"; \
			))