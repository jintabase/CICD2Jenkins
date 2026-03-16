package logic

import (
	"context"
	"testing"
	"time"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
)

type stubArticleRepository struct {
	articles map[string]model.Article
}

func (s *stubArticleRepository) Create(_ context.Context, article model.Article) (*model.Article, error) {
	if s.articles == nil {
		s.articles = map[string]model.Article{}
	}
	s.articles[article.ID] = article
	return &article, nil
}

func (s *stubArticleRepository) Update(_ context.Context, article model.Article) (*model.Article, error) {
	if s.articles == nil {
		s.articles = map[string]model.Article{}
	}
	s.articles[article.ID] = article
	return &article, nil
}

func (s *stubArticleRepository) Delete(_ context.Context, id string) error {
	if _, ok := s.articles[id]; !ok {
		return apperrors.ErrNotFound
	}
	delete(s.articles, id)
	return nil
}

func (s *stubArticleRepository) GetByID(_ context.Context, id string) (*model.Article, error) {
	article, ok := s.articles[id]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return &article, nil
}

func (s *stubArticleRepository) List(context.Context) ([]model.Article, error) {
	articles := make([]model.Article, 0, len(s.articles))
	for _, article := range s.articles {
		articles = append(articles, article)
	}
	return articles, nil
}

func TestArticleLogicCreateRequiresAdmin(t *testing.T) {
	svc := NewArticleLogic(&stubArticleRepository{})

	_, err := svc.Create(context.Background(), model.Actor{
		UserID:   "u-reader",
		Username: "reader",
		Role:     model.RoleUser,
	}, UpsertArticleInput{
		Title:   "hello",
		Content: "world",
	})
	if err != apperrors.ErrForbidden {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestArticleLogicCRUD(t *testing.T) {
	repo := &stubArticleRepository{
		articles: map[string]model.Article{
			"existing": {
				ID:         "existing",
				Title:      "old",
				Content:    "old content",
				AuthorID:   "admin-id",
				AuthorName: "admin",
				CreatedAt:  time.Now().UTC().Add(-time.Hour),
				UpdatedAt:  time.Now().UTC().Add(-time.Hour),
			},
		},
	}
	svc := NewArticleLogic(repo)
	admin := model.Actor{
		UserID:   "admin-id",
		Username: "admin",
		Role:     model.RoleSuperAdmin,
	}

	created, err := svc.Create(context.Background(), admin, UpsertArticleInput{
		Title:     "new article",
		Summary:   "summary",
		Content:   "content",
		Tags:      []string{"go", "Go", "blog"},
		Published: true,
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	if len(created.Tags) != 2 {
		t.Fatalf("expected deduplicated tags, got %v", created.Tags)
	}

	updated, err := svc.Update(context.Background(), admin, "existing", UpsertArticleInput{
		Title:     "updated",
		Summary:   "updated summary",
		Content:   "updated content",
		Tags:      []string{"updated"},
		Published: false,
	})
	if err != nil {
		t.Fatalf("update article: %v", err)
	}
	if updated.Title != "updated" {
		t.Fatalf("expected updated title, got %s", updated.Title)
	}

	if err := svc.Delete(context.Background(), admin, "existing"); err != nil {
		t.Fatalf("delete article: %v", err)
	}
}
