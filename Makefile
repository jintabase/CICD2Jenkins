APP_NAME=blog-api

.PHONY: run test fmt wire es-up es-down

run:
	go run ./cmd/blog-api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

wire:
	go run github.com/google/wire/cmd/wire ./internal/app

es-up:
	docker compose up -d elasticsearch

es-down:
	docker compose down
