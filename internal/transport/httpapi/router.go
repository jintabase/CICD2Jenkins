package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"cicd2jenkins/internal/domain"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi/handler"
	"cicd2jenkins/internal/transport/httpapi/httpx"
	authmiddleware "cicd2jenkins/internal/transport/httpapi/middleware"
)

func NewRouter(authService *service.AuthService, articleService *service.ArticleService) http.Handler {
	router := chi.NewRouter()
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(30 * time.Second))

	authHandler := handler.NewAuthHandler(authService)
	articleHandler := handler.NewArticleHandler(articleService)
	authmw := authmiddleware.NewAuthMiddleware(authService)

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authmw.Authenticate)

			r.Get("/me", authHandler.Me)

			r.Route("/articles", func(r chi.Router) {
				r.Get("/", articleHandler.List)
				r.Get("/{articleID}", articleHandler.GetByID)

				r.With(authmiddleware.RequireRoles(domain.RoleSuperAdmin)).Post("/", articleHandler.Create)
				r.With(authmiddleware.RequireRoles(domain.RoleSuperAdmin)).Put("/{articleID}", articleHandler.Update)
				r.With(authmiddleware.RequireRoles(domain.RoleSuperAdmin)).Delete("/{articleID}", articleHandler.Delete)
			})
		})
	})

	return router
}
