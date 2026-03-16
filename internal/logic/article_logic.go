package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
)

type articleRepository interface {
	Create(ctx context.Context, article model.Article) (*model.Article, error)
	Update(ctx context.Context, article model.Article) (*model.Article, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Article, error)
	List(ctx context.Context) ([]model.Article, error)
}

type ArticleLogic struct {
	repo articleRepository
}

type UpsertArticleInput struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

func NewArticleLogic(repo articleRepository) *ArticleLogic {
	return &ArticleLogic{repo: repo}
}

func (l *ArticleLogic) List(ctx context.Context) ([]model.Article, error) {
	articles, err := l.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	return articles, nil
}

func (l *ArticleLogic) GetByID(ctx context.Context, id string) (*model.Article, error) {
	article, err := l.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, fmt.Errorf("get article by id: %w", err)
	}
	return article, nil
}

func (l *ArticleLogic) Create(ctx context.Context, actor model.Actor, input UpsertArticleInput) (*model.Article, error) {
	if actor.Role != model.RoleSuperAdmin {
		return nil, apperrors.ErrForbidden
	}
	if err := validateArticleInput(input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	article := model.Article{
		ID:         uuid.NewString(),
		Title:      strings.TrimSpace(input.Title),
		Summary:    strings.TrimSpace(input.Summary),
		Content:    strings.TrimSpace(input.Content),
		Tags:       normalizeTags(input.Tags),
		Published:  input.Published,
		AuthorID:   actor.UserID,
		AuthorName: actor.Username,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	created, err := l.repo.Create(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("create article: %w", err)
	}
	return created, nil
}

func (l *ArticleLogic) Update(ctx context.Context, actor model.Actor, id string, input UpsertArticleInput) (*model.Article, error) {
	if actor.Role != model.RoleSuperAdmin {
		return nil, apperrors.ErrForbidden
	}
	if err := validateArticleInput(input); err != nil {
		return nil, err
	}

	current, err := l.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, fmt.Errorf("find existing article: %w", err)
	}

	current.Title = strings.TrimSpace(input.Title)
	current.Summary = strings.TrimSpace(input.Summary)
	current.Content = strings.TrimSpace(input.Content)
	current.Tags = normalizeTags(input.Tags)
	current.Published = input.Published
	current.UpdatedAt = time.Now().UTC()

	updated, err := l.repo.Update(ctx, *current)
	if err != nil {
		return nil, fmt.Errorf("update article: %w", err)
	}
	return updated, nil
}

func (l *ArticleLogic) Delete(ctx context.Context, actor model.Actor, id string) error {
	if actor.Role != model.RoleSuperAdmin {
		return apperrors.ErrForbidden
	}
	if strings.TrimSpace(id) == "" {
		return apperrors.ErrBadRequest
	}
	if err := l.repo.Delete(ctx, strings.TrimSpace(id)); err != nil {
		return fmt.Errorf("delete article: %w", err)
	}
	return nil
}

func validateArticleInput(input UpsertArticleInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("%w: title is required", apperrors.ErrBadRequest)
	}
	if strings.TrimSpace(input.Content) == "" {
		return fmt.Errorf("%w: content is required", apperrors.ErrBadRequest)
	}
	return nil
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return normalized
}
