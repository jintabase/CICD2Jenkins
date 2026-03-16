package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(ctx context.Context, article model.Article) (*model.Article, error) {
	if err := r.db.WithContext(ctx).Create(&article).Error; err != nil {
		return nil, fmt.Errorf("create article: %w", err)
	}
	return &article, nil
}

func (r *ArticleRepository) Update(ctx context.Context, article model.Article) (*model.Article, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", strings.TrimSpace(article.ID)).
		Updates(map[string]any{
			"title":       article.Title,
			"summary":     article.Summary,
			"content":     article.Content,
			"tags":        article.Tags,
			"published":   article.Published,
			"author_id":   article.AuthorID,
			"author_name": article.AuthorName,
			"updated_at":  article.UpdatedAt,
		})
	if result.Error != nil {
		return nil, fmt.Errorf("update article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}
	return &article, nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Article{}, "id = ?", strings.TrimSpace(id))
	if result.Error != nil {
		return fmt.Errorf("delete article: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *ArticleRepository) GetByID(ctx context.Context, id string) (*model.Article, error) {
	var article model.Article
	if err := r.db.WithContext(ctx).First(&article, "id = ?", strings.TrimSpace(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("get article by id: %w", err)
	}
	return &article, nil
}

func (r *ArticleRepository) List(ctx context.Context) ([]model.Article, error) {
	var articles []model.Article
	if err := r.db.WithContext(ctx).Order("updated_at DESC").Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	return articles, nil
}
