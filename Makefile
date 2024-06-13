.PHONY: run test help default

default: test

help:
	@echo 'Available commands:'
	@echo
	@echo 'Usage:'
	@echo '    make run           Compile the project.'
	@echo '    make test          Run tests on a compiled project.'
	@echo '    make format        Source code linting.'
	@echo

run:
	@go run ./cmd/server/main.go

test:
	go test -race ./...

format:
	go mod tidy
	go vet ./...
	go fmt ./...
	goimports -w -local github.com/wojciechpawlinow .