SHELL:=/bin/sh
.PHONY: all build test clean

export GO111MODULE=on

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}/build/bin

# Version
RELEASE_VER := $(shell git tag --list --sort=-creatordate  "v*" | head -n 1 )

# Go MOD
GO_MOD := $(shell go list -m)

# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif


# go source files, ignore vendor directory
SOURCE = $(shell find ${MKFILE_DIR} -type f -name "*.go")
TARGET = ${RELEASE_DIR}/guardianprobe

all: ${TARGET}

${TARGET}: ${SOURCE}
	mkdir -p ${RELEASE_DIR}
	go mod tidy
	CGO_ENABLED=0 go build -a -ldflags "-s -w -extldflags -static -X ${GO_MOD}/global.Ver=${RELEASE_VER}" -o ${TARGET} ${GO_MOD}/cmd/guardianprobe

build: all

test:
	go test -gcflags=-l -cover -race ${TEST_FLAGS} -v ./...

docker:
	DOCKER_BUILDKIT=1 docker build -t vmo2apps.azurecr.io/guardian/guardianprobe -f ${MKFILE_DIR}/resources/Dockerfile ${MKFILE_DIR}

release-docker:
	DOCKER_BUILDKIT=1 docker build -t vmo2apps.azurecr.io/guardian/guardianprobe:latest -f ${MKFILE_DIR}/resources/Dockerfile ${MKFILE_DIR}
	az account set --subscription O2UK-IT-CoreCommonServices-Tier0
	az acr login --name vmo2apps
	docker push vmo2apps.azurecr.io/guardian/guardianprobe:latest

clean:
	@rm -rf ${MKFILE_DIR}/build
