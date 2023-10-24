
VERSION ?= $(shell cat ./VERSION)
GO := docker run --rm -w /app -v $$PWD:/app -e CGO_ENABLED=0 golang:1.21-bookworm go

.PHONY: build

build:
	mkdir -p bin/
	$(GO) build -o ./bin/jobser-$(VERSION) -mod=vendor  -ldflags="-X 'main.Version=$(VERSION)'" ./cmd/jobserweb/main.go

test:
	$(GO) test -mod=vendor ./...
