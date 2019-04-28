SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

get-deps:
	dep ensure

cp-config:
	cp config_example.toml config.toml

build:
	go fmt ./...
	go build

clean:
	rm -fr data/

.PHONY: cp-config build
