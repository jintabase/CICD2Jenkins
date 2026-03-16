package repository

import (
	"context"

	"cicd2jenkins/internal/domain"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type ArticleRepository interface {
	EnsureIndex(ctx context.Context) error
	Create(ctx context.Context, article domain.Article) (*domain.Article, error)
	Update(ctx context.Context, article domain.Article) (*domain.Article, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	List(ctx context.Context) ([]domain.Article, error)
}
