package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
	"cicd2jenkins/internal/repository"
)

type ArticleService struct {
	repo repository.ArticleRepository
}

type UpsertArticleInput struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

func NewArticleService(repo repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) List(ctx context.Context) ([]domain.Article, error) {
	articles, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	return articles, nil
}

func (s *ArticleService) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	article, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, fmt.Errorf("get article by id: %w", err)
	}
	return article, nil
}

func (s *ArticleService) Create(ctx context.Context, actor domain.Actor, input UpsertArticleInput) (*domain.Article, error) {
	if actor.Role != domain.RoleSuperAdmin {
		return nil, apperrors.ErrForbidden
	}

	if err := validateArticleInput(input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	article := domain.Article{
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

	created, err := s.repo.Create(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("create article: %w", err)
	}
	return created, nil
}

func (s *ArticleService) Update(ctx context.Context, actor domain.Actor, id string, input UpsertArticleInput) (*domain.Article, error) {
	if actor.Role != domain.RoleSuperAdmin {
		return nil, apperrors.ErrForbidden
	}

	if err := validateArticleInput(input); err != nil {
		return nil, err
	}

	current, err := s.repo.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, fmt.Errorf("find existing article: %w", err)
	}

	current.Title = strings.TrimSpace(input.Title)
	current.Summary = strings.TrimSpace(input.Summary)
	current.Content = strings.TrimSpace(input.Content)
	current.Tags = normalizeTags(input.Tags)
	current.Published = input.Published
	current.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, *current)
	if err != nil {
		return nil, fmt.Errorf("update article: %w", err)
	}
	return updated, nil
}

func (s *ArticleService) Delete(ctx context.Context, actor domain.Actor, id string) error {
	if actor.Role != domain.RoleSuperAdmin {
		return apperrors.ErrForbidden
	}

	if strings.TrimSpace(id) == "" {
		return apperrors.ErrBadRequest
	}

	if err := s.repo.Delete(ctx, strings.TrimSpace(id)); err != nil {
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
