.PHONY: dev build test lint docker-up docker-down seed

dev:
	go run ./cmd/api

build:
	go build ./...

test:
	go test ./...

lint:
	golangci-lint run || true

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v
