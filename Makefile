default: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s %s\n\033[0m", $$1, $$2}'

# general
APP_NAME            = hostdb-collector-aws
APP_VER             = 0.1
CIRCLE_WORKFLOW_ID ?= ""
WORK_DIR            = $(shell pwd)

# git
ifeq ($(CIRCLECI), true)
GIT_BRANCH_DIRTY = $(CIRCLE_BRANCH)
else
GIT_BRANCH_DIRTY = $(shell git rev-parse --abbrev-ref HEAD)
endif
GIT_BRANCH     = $(shell echo "$(GIT_BRANCH_DIRTY)" | sed s/[[:punct:]]/_/g | tr '[:upper:]' '[:lower:]')
GIT_COMMIT_MSG = $(shell git log -1 --pretty=%B | tr '[:cntrl:]' ',' | sed 's/,*$$//g' | sed 's/,/, /g')
GIT_COMMIT_NUM = $(shell git rev-list --count HEAD)
GIT_COMMIT_SHA = $(shell git describe --tags --match '[0-9]*' --dirty --always --long)

# version
ifeq ($(GIT_BRANCH), master)
TAG     = latest
VERSION = $(APP_VER).$(GIT_COMMIT_NUM)
else
TAG     = $(GIT_BRANCH)
VERSION = $(APP_VER).$(GIT_COMMIT_NUM)-$(GIT_BRANCH)
endif
export TAG
export VERSION

# container
CONTAINER_REPO       = registry.pdxfixit.com
CONTAINER_IMAGE_NAME = $(CONTAINER_REPO)/$(APP_NAME)

# docker env file
ifneq ("$(wildcard env.list)","")
	DOCKER_ENV=--env-file ./env.list
endif
DOCKER_RUN_OPTIONS=-it --rm --name $(APP_NAME) $(DOCKER_ENV)

# use a go mod proxy so that we don't need deploy keys for GHE
export GOPROXY = https://goproxy.pdxfixit.com
# disable go mod sum checking for now
export GONOSUMDB = *.pdxfixit.com,github.com/pdxfixit/*

.PHONY: all
all: get test compile build push

.PHONY: get
get: ## go get will ensure dependencies are present
	go get -d

.PHONY: fmt
fmt: ## go fmt
	go fmt ./...

.PHONY: vet
vet: ## go vet
	go vet -v ./...

.PHONY: lint
lint: ## golint
ifeq (, $(shell which golint))
	go install golang.org/x/lint/golint@latest
endif
	go list ./... | xargs -L1 golint

.PHONY: errcheck
errcheck: ## errcheck
ifeq (, $(shell which errcheck))
	go install github.com/kisielk/errcheck@latest
endif
	errcheck ./...

.PHONY: test
test: fmt vet lint errcheck ## run the golang tests
	go test -v --failfast

.PHONY: compile
compile: $(APP_NAME) ## compile the linux/amd64 binary with c lib bindings for use in a scratch container

$(APP_NAME):
ifeq (, $(shell which gox))
	go install github.com/mitchellh/gox@latest
endif
	env CGO_ENABLED=0 gox -osarch="linux/amd64" -tags netgo -output $(APP_NAME)

.PHONY: build
build: ## create container image
	if [ ! -f $(APP_NAME) ]; then $(MAKE) compile; fi
	docker build -t $(APP_NAME) --label "version=$(VERSION)" .

.PHONY: run
run: ## run the container
	if [ "$$(docker images -q $(APP_NAME))" == "" ]; then $(MAKE) build; fi
	docker run $(DOCKER_RUN_OPTIONS) $(APP_NAME)

.PHONY: push
push: ## push container image to registry
ifeq ($(strip $(REGISTRY_USER)),)
	$(error "Username required (e.g. make push REGISTRY_USER=username REGISTRY_PASS=password)")
endif
ifeq ($(strip $(REGISTRY_PASS)),)
	$(error "Password required (e.g. make push REGISTRY_USER=username REGISTRY_PASS=password)")
endif
ifndef NEWRELIC_APIKEY
	$(error "New Relic API Key required (e.g. make push NEWRELIC_APIKEY=abc123)")
endif
	if [ "$$(docker images -q $(APP_NAME))" == "" ]; then $(MAKE) build; fi
	docker tag $(APP_NAME) $(CONTAINER_IMAGE_NAME):$(GIT_COMMIT_SHA)
	docker tag $(APP_NAME) $(CONTAINER_IMAGE_NAME):$(VERSION)
	docker tag $(APP_NAME) $(CONTAINER_IMAGE_NAME):$(TAG)

	@echo $(REGISTRY_PASS) | docker login -u $(REGISTRY_USER) --password-stdin $(CONTAINER_REPO)
	docker push $(CONTAINER_IMAGE_NAME):$(GIT_COMMIT_SHA)
	docker push $(CONTAINER_IMAGE_NAME):$(VERSION)
	docker push $(CONTAINER_IMAGE_NAME):$(TAG)

	jq -n --arg msg "$(GIT_COMMIT_MSG)" '{"deployment":{"revision":"$(GIT_COMMIT_SHA)","description":$$msg,"user":"$(APP_NAME)"}}' > json.txt
	curl -i -X POST -H 'X-Api-Key: $(NEWRELIC_APIKEY)' -H 'Content-Type: application/json' -d @json.txt 'https://api.newrelic.com/v2/applications/1234567890/deployments.json'
	rm json.txt

.PHONY: sample-data
sample-data: ## generate HostDB JSON data files in ./sample-data
#ifeq ($(strip $(AWS_ACCESS_KEY_ID)),)
#	$(error "AWS Access Key is required (e.g. make sample-data AWS_ACCESS_KEY_ID=accesskeyid AWS_SECRET_ACCESS_KEY=sekret)")
#endif
#ifeq ($(strip $(AWS_SECRET_ACCESS_KEY)),)
#	$(error "AWS Secret Access Key is required (e.g. make sample-data AWS_ACCESS_KEY_ID=accesskeyid AWS_SECRET_ACCESS_KEY=sekret)")
#endif
	if [ "$$(docker images -q $(SOURCE_TAG))" == "" ]; then $(MAKE) build; fi
	mkdir -p sample-data
	docker run $(DOCKER_RUN_OPTIONS) -e HOSTDB_COLLECTOR_AWS_COLLECTOR_SAMPLE_DATA=true -v `pwd`/sample-data:/sample-data $(SOURCE_TAG)

.PHONY: stop
stop: ## stop the container
	if [ "$$(docker ps -a -q -f 'name=$(APP_NAME)')" ]; then docker stop -t0 $(APP_NAME); fi

.PHONY: clean
clean: ## clean up any artifacts
	rm -f $(WORK_DIR)/$(APP_NAME)
	rm -rf $(WORK_DIR)/sample-data
	if [ "$$(docker ps -a -q -f 'name=$(APP_NAME)')" ]; then docker stop -t0 $(APP_NAME); fi
	if [ "$$(docker images -q $(APP_NAME))" ]; then docker rmi $(APP_NAME); fi
