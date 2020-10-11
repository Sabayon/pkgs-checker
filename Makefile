NAME := pkgs-checker
PACKAGE_NAME ?= $(NAME)
REVISION := $(shell git rev-parse --short HEAD || echo unknown)
BUILD_PLATFORMS ?= -osarch="linux/amd64" -osarch="linux/386" -osarch="linux/arm"

.PHONY: all

all: pkgs-checker

.PHONY: pkgs-checker
pkgs-checker:
	CGO_ENABLE=0 go build -v .

.PHONY: test
test:
	go test -v -tags all -cover -race ./...
	#go test -v -tags all -cover -race ./... -ginkgo.v

.PHONY: clean
clean:
	-rm pkgs-checker
	-rm -rf release/

.PHONY: deps
deps:
	go env
	# Installing dependencies...
	go get golang.org/x/lint/golint
	go get github.com/mitchellh/gox
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/onsi/gomega/...

.PHONY: multiarch-build-dev
multiarch-build-dev: deps
	CGO_ENABLE=0 gox $(BUILD_PLATFORMS) -output="release/$(NAME)-$(REVISION)-{{.OS}}-{{.Arch}}" -ldflags "-extldflags=-Wl,--allow-multiple-definition"
