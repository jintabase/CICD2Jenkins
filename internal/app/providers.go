package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"cicd2jenkins/internal/config"
	"cicd2jenkins/internal/domain"
	"cicd2jenkins/internal/repository"
	"cicd2jenkins/internal/repository/elasticsearch"
	"cicd2jenkins/internal/repository/memory"
	"cicd2jenkins/internal/service"
)

func provideElasticsearchClient(cfg config.Config) (*es8.Client, error) {
	client, err := es8.NewClient(es8.Config{
		Addresses: cfg.Elasticsearch.Addresses,
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
		Transport: &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	return client, nil
}

func provideArticleRepository(cfg config.Config, client *es8.Client) (*elasticsearch.ArticleRepository, error) {
	articleRepo := elasticsearch.NewArticleRepository(client, cfg.Elasticsearch.Index)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Elasticsearch.RequestTimeout)
	defer cancel()

	if err := articleRepo.EnsureIndex(ctx); err != nil {
		return nil, fmt.Errorf("prepare elasticsearch index: %w", err)
	}

	return articleRepo, nil
}

func provideSeedUsers(cfg config.Config) ([]domain.User, error) {
	users := make([]domain.User, 0, len(cfg.SeedUsers))
	now := time.Now().UTC()

	for _, seed := range cfg.SeedUsers {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash seed user password: %w", err)
		}

		users = append(users, domain.User{
			ID:           uuid.NewString(),
			Username:     seed.Username,
			PasswordHash: string(passwordHash),
			Role:         seed.Role,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	return users, nil
}

func provideUserRepository(users []domain.User) *memory.UserRepository {
	return memory.NewUserRepository(users)
}

func provideAuthService(cfg config.Config, users repository.UserRepository) *service.AuthService {
	return service.NewAuthService(users, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
}

func provideHTTPServer(cfg config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}
