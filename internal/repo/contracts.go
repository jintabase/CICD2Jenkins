package repo

import (
	"context"

	"cicd2jenkins/internal/model"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type ArticleRepository interface {
	Create(ctx context.Context, article model.Article) (*model.Article, error)
	Update(ctx context.Context, article model.Article) (*model.Article, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	List(ctx context.Context) ([]model.Article, error)
}
