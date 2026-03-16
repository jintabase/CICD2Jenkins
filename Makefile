APP_NAME=blog-api

.PHONY: run test fmt wire mysql-up mysql-down

run:
	go run ./cmd/blog-api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

wire:
	go run github.com/google/wire/cmd/wire ./internal/app

mysql-up:
	docker compose up -d mysql

mysql-down:
	docker compose down
