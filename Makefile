.PHONY: build
build:
	go build -o ./share_bot ./cmd

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

migrate-up:
	goose -dir ./migrations sqlite3 ./data.db up

.DEFAULT_GOAL := build