APP := bruce
APP_ENTRY := cmd/main.go
SHELL := /bin/bash
VER := $(shell git rev-parse --short HEAD)
ifndef VER
$(error VER not set: Run make this way: `make VER=1.0.31`)
endif
ROOTPATH := $(shell echo ${PWD})
ARCH := $(shell uname -m)
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')

.PHONY: all clean setup build deploy run

all: build

clean:
	rm -rf ${ROOTPATH}/.build/
	rm -rf ${ROOTPATH}/vendor/

setup:
	mkdir -p ${ROOTPATH}/.build/bin

package: build zipit

deploy: build-local push

build: clean setup build-local 

build-local:
	@go version
	@go get ./cmd/...
	GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w -X main.version=${VER}" -o ${ROOTPATH}/.build/bin/${APP}-${OS}-${ARCH} ${APP_ENTRY}


zipit:
	cd .build/ && zip -r ${APP}-${VER}-${OS}-${ARCH}.zip ./*
	@echo "package ready under: .build/"
