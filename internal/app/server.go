package app

import (
	"context"
	"net/http"

	"cicd2jenkins/internal/config"
)

//go:generate go run github.com/google/wire/cmd/wire ./internal/app

func NewServer(cfg config.Config) (*http.Server, func(context.Context) error, error) {
	server, cleanup, err := initializeServer(cfg)
	if err != nil {
		return nil, nil, err
	}

	if cleanup == nil {
		return server, func(context.Context) error { return nil }, nil
	}

	return server, func(context.Context) error {
		cleanup()
		return nil
	}, nil
}
