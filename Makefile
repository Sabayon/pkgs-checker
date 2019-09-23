NAME := pkgs-checker
PACKAGE_NAME ?= $(NAME)
REVISION := $(shell git rev-parse --short HEAD || echo unknown)
BUILD_PLATFORMS ?= -osarch="linux/amd64" -osarch="linux/386" -osarch="linux/arm"

.PHONY: all

all: pkgs-checker

.PHONY: pkgs-checker
pkgs-checker:
	go build -v .

.PHONY: test
test:
	go test -v -tags all -cover -race ./...
	#go test -v -tags all -cover -race ./... -ginkgo.v

.PHONY: clean
clean:
	-rm pkgs-checker
	-rm -rf release/

.PHONY: multiarch-build-dev
multiarch-build-dev:
	gox $(BUILD_PLATFORMS) -output="release/$(NAME)-$(REVISION)-{{.OS}}-{{.Arch}}" -ldflags "-extldflags=-Wl,--allow-multiple-definition"


