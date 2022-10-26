.PHONY: build
build:
	CGO_ENABLED=0 go build -v ./cmd/share_bot

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
