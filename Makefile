.DELETE_ON_ERROR: clean

EXECUTABLES = go zip shasum
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

PROJECT_NAME      ?= idp-scim-sync
PROJECT_NAMESPACE ?= slashdevops
PROJECT_MODULES_PATH := $(shell ls -d cmd/*)
PROJECT_MODULES_NAME := $(foreach dir_name, $(PROJECT_MODULES_PATH), $(shell basename $(dir_name)) )
PROJECT_DEPENDENCIES := $(shell go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all)

BUILD_DIR       := ./build
DIST_DIR        := ./dist
DIST_ASSEST_DIR := $(DIST_DIR)/assets

PROJECT_COVERAGE_FILE ?= $(BUILD_DIR)/coverage.txt
PROJECT_COVERAGE_MODE	?= atomic
PROJECT_COVERAGE_TAGS ?= unit

GIT_VERSION  ?= $(shell git rev-parse --abbrev-ref HEAD | cut -d "/" -f 2)
GIT_REVISION ?= $(shell git rev-parse HEAD | tr -d '\040\011\012\015\n')
GIT_BRANCH   ?= $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')
GIT_USER     ?= $(shell git config --get user.name | tr -d '\040\011\012\015\n')
BUILD_DATE   ?= $(shell date +'%Y-%m-%dT%H:%M:%S')

GO_LDFLAGS_OPTIONS ?= -s -w
define EXTRA_GO_LDFLAGS_OPTIONS
-X '"'github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Version=$(GIT_VERSION)'"' \
-X '"'github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Revision=$(GIT_REVISION)'"' \
-X '"'github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.Branch=$(GIT_BRANCH)'"' \
-X '"'github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.BuildUser=$(GIT_USER)'"' \
-X '"'github.com/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)/internal/version.BuildDate=$(BUILD_DATE)'"'
endef
GO_LDFLAGS     := -ldflags "$(GO_LDFLAGS_OPTIONS) $(EXTRA_GO_LDFLAGS_OPTIONS)"
GO_CGO_ENABLED ?= 0
GO_OPTS        ?= -v
GO_OS          ?= darwin linux windows
GO_ARCH        ?= arm64 amd64
# avoid mocks in tests
GO_FILES       := $(shell go list ./... | grep -v /mocks/)
GO_GRAPH_FILE  := $(BUILD_DIR)/go-mod-graph.txt

CONTAINER_OS   ?= linux
CONTAINER_ARCH ?= arm64v8 amd64
CONTAINER_NAMESPACE ?= $(PROJECT_NAMESPACE)
CONTAINER_IMAGE_NAME ?= $(PROJECT_NAME)

DOCKER_CONTAINER_REPO  ?= docker.io
GITHUB_CONTAINER_REPO  ?= ghcr.io
AWS_ECR_CONTAINER_REPO ?= public.ecr.aws/l2n7y5s7

AWS_SAM_LAMBDA_BINARY_NAME ?= bootstrap
AWS_SAM_PROJECT_APP_NAME   ?= idpscim
AWS_SAM_OS                 ?= linux
AWS_SAM_ARCH               ?= arm64

######## Functions ########
# this is a function that will execute a command and print a message
# MAKE_DEBUG=true make <target> will print the command
# MAKE_STOP_ON_ERRORS=true make any fail will stop the execution if the command fails, this is useful for CI
# NOTE: if the command has a > it will print the output into the original redirect of the command
MAKE_STOP_ON_ERRORS := false
MAKE_DEBUG := false

define exec_cmd
$(if $(filter $(MAKE_DEBUG),true),\
	${1} \
, \
	$(if $(filter $(MAKE_STOP_ON_ERRORS),true),\
		@ERROR_OCCURRED=0; ${1} > /dev/null || ERROR_OCCURRED=1; if [ $$ERROR_OCCURRED -eq 0 ]; then printf "  🤞 ${1} ✅\n"; else printf "  ${1} ❌ 🖕\n"; exit 1; fi \
	, \
		$(if $(findstring >, $1),\
			@${1} 2>/dev/null && printf "  🤞 ${1} ✅\n" || printf "  ${1} ❌ 🖕\n" \
		, \
			@${1} > /dev/null 2>&1 && printf '  🤞 ${1} ✅\n' || printf '  ${1} ❌ 🖕\n' \
		) \
	) \
)

endef # don't remove the white space at the end of the line
# this is a function that will execute a command and print a message

###############################################################################
######## Targets ##############################################################
##@ Default command
.PHONY: all
all: clean test build ## Clean, test and build the application.  Execute by default when make is called without arguments

###############################################################################
##@ Golang commands
.PHONY: go-fmt
go-fmt: ## Format go code
	@printf "👉 Formatting go code...\n"
	$(call exec_cmd, go fmt ./... )

.PHONY: go-vet
go-vet: ## Vet go code
	@printf "👉 Vet go code...\n"
	$(call exec_cmd, go vet ./... )

.PHONY: go-generate
go-generate: ## Generate go code
	@printf "👉 Generating go code...\n"
	$(call exec_cmd, go generate ./... )

.PHONY: go-mod-tidy
go-mod-tidy: ## Clean go.mod and go.sum
	@printf "👉 Cleaning go.mod and go.sum...\n"
	$(call exec_cmd, go mod tidy)

.PHONY: go-mod-update
go-mod-update: go-mod-tidy ## Update go.mod and go.sum
	@printf "👉 Updating go.mod and go.sum...\n"
	$(foreach DEP, $(PROJECT_DEPENDENCIES), \
		$(call exec_cmd, go get -u $(DEP)) \
	)

.PHONY: go-mod-vendor
go-mod-vendor: ## Create mod vendor
	@printf "👉 Creating mod vendor...\n"
	$(call exec_cmd, go mod vendor)

.PHONY: go-mod-verify
go-mod-verify: ## Verify go.mod and go.sum
	@printf "👉 Verifying go.mod and go.sum...\n"
	$(call exec_cmd, go mod verify)

.PHONY: go-mod-download
go-mod-download: ## Download go dependencies
	@printf "👉 Downloading go dependencies...\n"
	$(call exec_cmd, go mod download)

.PHONY: go-mod-graph
go-mod-graph: ## Create a file with the go dependencies graph in build dir
	@printf "👉 Printing go dependencies graph...\n"
	$(call exec_cmd, go mod graph > $(GO_GRAPH_FILE))

# this target is needed to create the dist folder and the coverage file
$(PROJECT_COVERAGE_FILE):
	@printf "👉 Creating coverage file...\n"
	$(call exec_cmd, mkdir -p $(BUILD_DIR) )
	$(call exec_cmd, touch $(PROJECT_COVERAGE_FILE) )

.PHONY: go-test-coverage
go-test-coverage: test ## Shows in you browser the test coverage report per package
	@printf "👉 Running got tool coverage...\n"
	$(call exec_cmd, go tool cover -html=$(PROJECT_COVERAGE_FILE))

###############################################################################
##@ Test commands
.PHONY: test
test: $(PROJECT_COVERAGE_FILE) go-mod-tidy go-fmt go-vet go-generate ## Run tests
	@printf "👉 Running tests...\n"
	$(call exec_cmd, go test \
		-v -race \
		-coverprofile=$(PROJECT_COVERAGE_FILE) \
		-covermode=$(PROJECT_COVERAGE_MODE) \
			-tags=$(PROJECT_COVERAGE_TAGS) \
		./... \
	)

###############################################################################
##@ Build commands
.PHONY: build
build: go-generate go-fmt go-vet test ## Build the application
	@printf "👉 Building applications...\n"
	$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
		$(call exec_cmd, CGO_ENABLED=$(GO_CGO_ENABLED) go build $(GO_LDFLAGS) $(GO_OPTS) -o $(BUILD_DIR)/$(proj_mod) ./cmd/$(proj_mod)/ ) \
		$(call exec_cmd, chmod +x $(BUILD_DIR)/$(proj_mod) ) \
	)

.PHONY: build-dist
build-dist: ## Build the application for all platforms defined in GO_OS and GO_ARCH in this Makefile
	@printf "👉 Building application for different platforms...\n"
	$(foreach GOOS, $(GO_OS), \
		$(foreach GOARCH, $(GO_ARCH), \
			$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
				$(call exec_cmd, GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(GO_CGO_ENABLED) go build $(GO_LDFLAGS) $(GO_OPTS) -o $(DIST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH) ./cmd/$(proj_mod)/ ) \
				$(call exec_cmd, chmod +x $(DIST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH)) \
			)\
		)\
	)

.PHONY: build-dist-zip
build-dist-zip: ## Build the application for all platforms defined in GO_OS and GO_ARCH in this Makefile and create a zip file for each binary
	@printf "👉 Creating zip files for distribution...\n"
	$(call exec_cmd, mkdir -p $(DIST_ASSEST_DIR))
	$(foreach GOOS, $(GO_OS), \
		$(foreach GOARCH, $(GO_ARCH), \
			$(foreach proj_mod, $(PROJECT_MODULES_NAME), \
				$(call exec_cmd, zip --junk-paths -r $(DIST_ASSEST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH).zip $(DIST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH) ) \
				$(call exec_cmd, shasum -a 256 $(DIST_ASSEST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH).zip | cut -d ' ' -f 1 > $(DIST_ASSEST_DIR)/$(proj_mod)-$(GOOS)-$(GOARCH).sha256 ) \
			) \
		) \
	)

###############################################################################
# This target is used by AWS SAM build command
# and was added to build the binary using custom flags
# Ref:
# + https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/building-custom-runtimes.html
# + https://jiangsc.me/2021/01/24/Lessons-Learnt-On-Deploying-GO-Lambda-Application-on-AWS/
# NOTES:
# + The ARTIFACTS_DIR environment variable is injected by AWS SAM build command
##@ AWS Lambda commands
.PHONY: build-LambdaFunction
build-LambdaFunction: ## Build the application for AWS Lambda, this target is used by AWS SAM build command
	@printf "👉 Called from sam build command ...\n"
	@printf "  👉 ARTIFACTS_DIR injected from sam build command: %s\n" $(ARTIFACTS_DIR)
	$(call exec_cmd, GOOS=$(AWS_SAM_OS) GOARCH=$(AWS_SAM_ARCH) CGO_ENABLED=$(GO_CGO_ENABLED) go build $(GO_LDFLAGS) $(GO_OPTS) -tags lambda.norpc -o $(ARTIFACTS_DIR)/$(AWS_SAM_LAMBDA_BINARY_NAME) ./cmd/$(AWS_SAM_PROJECT_APP_NAME)/ )

###############################################################################
##@ Container commands
.PHONY: container-build
container-build: build-dist ## Build the container image
	@printf "👉 Building container image...\n"
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring v, $(ARCH)), $(eval BIN_ARCH = arm64), $(eval BIN_ARCH = $(ARCH)) ) \
			$(call exec_cmd, docker build \
													--build-arg ARCH=$(ARCH) \
													--build-arg BIN_ARCH=$(BIN_ARCH) \
													--build-arg OS=$(OS) \
													--build-arg PROJECT_NAME=$(AWS_SAM_PROJECT_APP_NAME) \
													--tag $(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH) \
													--tag $(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH) \
													--tag $(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH) \
													--tag $(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH) \
													--tag $(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH) \
													--tag $(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH) \
													./. \
			) \
		) \
	)

.PHONY: container-publish-docker
container-publish-docker: ## Publish the container image to docker hub
	@printf "👉 Publishing container image to docker hub...\n"
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring v, $(ARCH)), $(eval BIN_ARCH = arm64), $(eval BIN_ARCH = $(ARCH)) ) \
			\
			$(call exec_cmd, docker push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" ) \
			$(call exec_cmd, docker push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" ) \
			\
			$(call exec_cmd, docker manifest create --amend \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
									--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest create --amend \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
									"$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
									--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" ) \
			$(call exec_cmd, docker manifest push "$(DOCKER_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" ) \
		) \
	)

.PHONY: container-publish-github
container-publish-github: ## Publish the container image to github container registry
	@printf "👉 Publishing container image to github container registry...\n"
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring v, $(ARCH)), $(eval BIN_ARCH = arm64), $(eval BIN_ARCH = $(ARCH)) ) \
			\
			$(call exec_cmd, docker push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" ) \
			$(call exec_cmd, docker push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" ) \
			\
			$(call exec_cmd, docker manifest create --amend \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
										--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest create --amend \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
										"$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
										--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" ) \
			$(call exec_cmd, docker manifest push "$(GITHUB_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" ) \
		) \
	)

.PHONY: container-publish-aws-ecr
container-publish-aws-ecr: ## Publish the container image to AWS ECR
	@printf "👉 Publishing container image to AWS ECR...\n"
	$(foreach OS, $(CONTAINER_OS), \
		$(foreach ARCH, $(CONTAINER_ARCH), \
			$(if $(findstring v, $(ARCH)), $(eval BIN_ARCH = arm64), $(eval BIN_ARCH = $(ARCH)) ) \
			\
			$(call exec_cmd, docker push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" ) \
			$(call exec_cmd, docker push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" ) \
			\
			$(call exec_cmd, docker manifest create --amend \
										"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
										"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
										"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" \
										"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)-$(OS)-$(ARCH)" \
										--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest create --amend \
				"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
				"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
				) \
			$(call exec_cmd, docker manifest annotate \
				"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" \
				"$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest-$(OS)-$(ARCH)" \
				--os $(OS) --arch $(BIN_ARCH) $(if $(findstring v, $(ARCH)), --variant "v8", ) \
				) \
			\
			$(call exec_cmd, docker manifest push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):latest" ) \
			$(call exec_cmd, docker manifest push "$(AWS_ECR_CONTAINER_REPO)/$(CONTAINER_NAMESPACE)/$(CONTAINER_IMAGE_NAME):$(GIT_VERSION)" ) \
		) \
	)

###############################################################################
##@ Support Commands
.PHONY: clean
clean: ## Clean the environment
	@printf "👉 Cleaning environment...\n"
	$(call exec_cmd, go clean -n -x -i)
	$(call exec_cmd, rm -rf $(BUILD_DIR) $(DIST_DIR) .aws-sam ./build.toml ./packaged.yaml )

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##";                                             \
		printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ \
		{ printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/            \
		{ printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } '                  \
		$(MAKEFILE_LIST)