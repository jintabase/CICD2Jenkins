package app

import (
	"context"
	"net/http"

	"cicd2jenkins/internal/config"
)

//go:generate go run github.com/google/wire/cmd/wire ./internal/app

func NewServer(cfg config.Config) (*http.Server, func(context.Context) error, error) {
	server, err := initializeServer(cfg)
	if err != nil {
		return nil, nil, err
	}

	return server, func(context.Context) error { return nil }, nil
}
