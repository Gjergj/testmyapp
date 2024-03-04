SHELL=/bin/bash -eo pipefail

VERSION = $(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $NF+=1; print $0}')
PWD = $(shell pwd)
GO ?= go

build:
	mkdir -p ./bin
	$(GO) build -o ./bin/ ./cmd/testmyapp/...

patch_release:
	$(eval VERSION=$(shell git describe --tags --abbrev=0 | awk -F. '{OFS="."; $$NF+=1; print $0}'))
	git tag -a $(VERSION) -m $(VERSION)
	git push origin $(VERSION)

