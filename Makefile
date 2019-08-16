.PHONY: all

all: pkgs-checker

.PHONY: pkgs-checker
pkgs-checker:
	go build -v .

clean:
	-rm pkgs-checker
