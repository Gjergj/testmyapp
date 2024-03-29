SHELL=/bin/bash -eo pipefail

VERSION = $(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $NF+=1; print $0}')
PWD = $(shell pwd)
GO ?= go
TAG := $(shell git describe --tags --abbrev=0)

build:
	mkdir -p ./bin
	$(GO) build -ldflags "-X main.version=$(TAG)" -o ./bin/ ./cmd/testmyapp/...

patch_release:
	$(eval VERSION=$(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $$NF+=1; print $0}'))
	git tag -a $(VERSION) -m $(VERSION)
	git push origin $(VERSION)

