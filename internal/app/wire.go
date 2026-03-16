//go:build wireinject

package app

import (
	"net/http"

	"github.com/google/wire"

	"cicd2jenkins/internal/config"
	"cicd2jenkins/internal/logic"
	"cicd2jenkins/internal/repo"
	"cicd2jenkins/internal/repo/gormrepo"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi"
	"cicd2jenkins/internal/transport/httpapi/middleware"
)

var serverSet = wire.NewSet(
	provideSeedUsers,
	provideDatabase,
	provideUserRepository,
	provideArticleRepository,
	provideAuthLogic,
	logic.NewArticleLogic,
	service.NewAuthService,
	service.NewArticleService,
	middleware.NewAuthMiddleware,
	httpapi.NewRouter,
	provideHTTPServer,
	wire.Bind(new(repo.ArticleRepository), new(*gormrepo.ArticleRepository)),
	wire.Bind(new(repo.UserRepository), new(*gormrepo.UserRepository)),
)

func initializeServer(cfg config.Config) (*http.Server, func(), error) {
	wire.Build(serverSet)
	return nil, nil, nil
}
