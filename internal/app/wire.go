//go:build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"

	"cicd2jenkins/internal/config"
	"cicd2jenkins/internal/repository"
	"cicd2jenkins/internal/repository/elasticsearch"
	"cicd2jenkins/internal/repository/memory"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi"
)

var serverSet = wire.NewSet(
	provideElasticsearchClient,
	provideArticleRepository,
	provideSeedUsers,
	provideUserRepository,
	provideAuthService,
	service.NewArticleService,
	httpapi.NewRouter,
	provideHTTPServer,
	wire.Bind(new(repository.ArticleRepository), new(*elasticsearch.ArticleRepository)),
	wire.Bind(new(repository.UserRepository), new(*memory.UserRepository)),
)

func initializeServer(cfg config.Config) (*http.Server, error) {
	wire.Build(serverSet)
	return nil, nil
}
