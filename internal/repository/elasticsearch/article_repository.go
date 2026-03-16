package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	es8 "github.com/elastic/go-elasticsearch/v8"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
)

type ArticleRepository struct {
	client *es8.Client
	index  string
}

func NewArticleRepository(client *es8.Client, index string) *ArticleRepository {
	return &ArticleRepository{
		client: client,
		index:  index,
	}
}

func (r *ArticleRepository) EnsureIndex(ctx context.Context) error {
	res, err := r.client.Indices.Exists([]string{r.index}, r.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("check index existence: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
	default:
		return fmt.Errorf("check index existence: %s", res.Status())
	}

	mapping := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"title":       map[string]any{"type": "text"},
				"summary":     map[string]any{"type": "text"},
				"content":     map[string]any{"type": "text"},
				"tags":        map[string]any{"type": "keyword"},
				"published":   map[string]any{"type": "boolean"},
				"author_id":   map[string]any{"type": "keyword"},
				"author_name": map[string]any{"type": "keyword"},
				"created_at":  map[string]any{"type": "date"},
				"updated_at":  map[string]any{"type": "date"},
			},
		},
	}

	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("marshal index mapping: %w", err)
	}

	createRes, err := r.client.Indices.Create(
		r.index,
		r.client.Indices.Create.WithContext(ctx),
		r.client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("create index: %s", readResponseMessage(createRes.Body))
	}

	return nil
}

func (r *ArticleRepository) Create(ctx context.Context, article domain.Article) (*domain.Article, error) {
	return r.indexDocument(ctx, article)
}

func (r *ArticleRepository) Update(ctx context.Context, article domain.Article) (*domain.Article, error) {
	return r.indexDocument(ctx, article)
}

func (r *ArticleRepository) Delete(ctx context.Context, id string) error {
	res, err := r.client.Delete(
		r.index,
		id,
		r.client.Delete.WithContext(ctx),
		r.client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("delete article: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return apperrors.ErrNotFound
	}
	if res.IsError() {
		return fmt.Errorf("delete article: %s", readResponseMessage(res.Body))
	}

	var payload struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode delete response: %w", err)
	}

	if payload.Result == "not_found" {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *ArticleRepository) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	res, err := r.client.Get(r.index, id, r.client.Get.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("get article: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, apperrors.ErrNotFound
	}
	if res.IsError() {
		return nil, fmt.Errorf("get article: %s", readResponseMessage(res.Body))
	}

	var payload struct {
		Found  bool           `json:"found"`
		ID     string         `json:"_id"`
		Source domain.Article `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode article: %w", err)
	}
	if !payload.Found {
		return nil, apperrors.ErrNotFound
	}

	payload.Source.ID = payload.ID
	return &payload.Source, nil
}

func (r *ArticleRepository) List(ctx context.Context) ([]domain.Article, error) {
	query := map[string]any{
		"size": 100,
		"sort": []map[string]any{
			{"updated_at": map[string]any{"order": "desc"}},
		},
		"query": map[string]any{
			"match_all": map[string]any{},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal search query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("list articles: %s", readResponseMessage(res.Body))
	}

	var payload struct {
		Hits struct {
			Hits []struct {
				ID     string         `json:"_id"`
				Source domain.Article `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	articles := make([]domain.Article, 0, len(payload.Hits.Hits))
	for _, hit := range payload.Hits.Hits {
		hit.Source.ID = hit.ID
		articles = append(articles, hit.Source)
	}

	return articles, nil
}

func (r *ArticleRepository) indexDocument(ctx context.Context, article domain.Article) (*domain.Article, error) {
	body, err := json.Marshal(article)
	if err != nil {
		return nil, fmt.Errorf("marshal article document: %w", err)
	}

	res, err := r.client.Index(
		r.index,
		bytes.NewReader(body),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithDocumentID(article.ID),
		r.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return nil, fmt.Errorf("index article: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("index article: %s", readResponseMessage(res.Body))
	}

	return &article, nil
}

func readResponseMessage(body io.Reader) string {
	payload, err := io.ReadAll(body)
	if err != nil {
		return "unknown elasticsearch error"
	}

	text := strings.TrimSpace(string(payload))
	if text == "" {
		return "unknown elasticsearch error"
	}
	return text
}
