.PHONY: all

all: pkgs-checker

.PHONY: pkgs-checker
pkgs-checker:
	go build -v .

.PHONY: test
test:
	go test -v -tags all -cover -race ./...
	#go test -v -tags all -cover -race ./... -ginkgo.v

clean:
	-rm pkgs-checker
