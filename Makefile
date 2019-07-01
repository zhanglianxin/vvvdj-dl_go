SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

GO_PACKAGE := github.com/zhanglianxin/vvvdj-dl_go
CROSS_TARGETS := linux/amd64 darwin/amd64 windows/386 windows/amd64

default: cp-config build cross gen-sha1

get-deps:
	dep ensure

cp-config:
	@-cp -n config_example.toml config.toml || echo 'config.toml exists'

build:
	go fmt ./...
	@#go build

clean:
	rm -fr data/*

cross:
	gox -osarch="$(CROSS_TARGETS)" $(GO_PACKAGE)
	@$(MAKE) gen-sha1

rm-sha1:
	@rm -f vvvdj-dl_*.sha1

gen-sha1: rm-sha1
	@$$(for f in $$(find vvvdj-dl_* -type f); do shasum $$f > $$f.sha1; done)
