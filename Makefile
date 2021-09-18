
NAME := pkgs-checker
PACKAGE_NAME ?= $(NAME)
REVISION := $(shell git rev-parse --short HEAD || echo dev)
VERSION := $(shell git describe --tags || echo $(REVISION))
VERSION := $(shell echo $(VERSION) | sed -e 's/^v//g')
BUILD_PLATFORMS ?= -osarch="linux/amd64" -osarch="linux/386" -osarch="linux/arm"

# go tool nm ./luet | grep Commit
override LDFLAGS += -X "github.com/Sabayon/pkgs-checker/cmd.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
override LDFLAGS += -X "github.com/Sabayon/pkgs-checker/cmd.BuildCommit=$(shell git rev-parse HEAD)"


.PHONY: all

all: pkgs-checker

.PHONY: pkgs-checker
pkgs-checker:
	# pkgs-checker uses go-sqlite3 that require CGO
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)'

.PHONY: test
test:
	go test -v -tags all -cover -race ./...
	#go test -v -tags all -cover -race ./... -ginkgo.v

.PHONY: coverage
coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: test-coverage
test-coverage:
	scripts/ginkgo.coverage.sh --codecov

.PHONY: clean
clean:
	-rm pkgs-checker
	-rm -rf release/ dist/

.PHONY: deps
deps:
	go env
	# Installing dependencies...
	GO111MODULE=off go get golang.org/x/lint/golint
	GO111MODULE=off go get github.com/mitchellh/gox
	GO111MODULE=off go get golang.org/x/tools/cmd/cover
	GO111MODULE=off go get github.com/onsi/ginkgo/ginkgo
	GO111MODULE=off go get github.com/onsi/gomega/...

.PHONY: goreleaser-snapshot
goreleaser-snapshot:
	rm -rf dist/ || true
	goreleaser release --debug --skip-publish  --skip-validate --snapshot

.PHONY: multiarch-build-dev
multiarch-build-dev: deps
	CGO_ENABLED=1 gox $(BUILD_PLATFORMS) -output="release/$(NAME)-$(REVISION)-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS) -extldflags=-Wl,--allow-multiple-definition"
