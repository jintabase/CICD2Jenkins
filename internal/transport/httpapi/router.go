package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cicd2jenkins/internal/model"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi/httpx"
	authmiddleware "cicd2jenkins/internal/transport/httpapi/middleware"
)

func NewRouter(authService *service.AuthService, articleService *service.ArticleService, authmw *authmiddleware.AuthMiddleware) http.Handler {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		httpx.WriteJSON(c, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	api.POST("/auth/login", authService.Login)

	protected := api.Group("/")
	protected.Use(authmw.Authenticate())
	protected.GET("/me", authService.Me)
	protected.GET("/articles", articleService.List)
	protected.GET("/articles/:articleID", articleService.GetByID)

	adminArticles := protected.Group("/articles")
	adminArticles.Use(authmiddleware.RequireRoles(model.RoleSuperAdmin))
	adminArticles.POST("", articleService.Create)
	adminArticles.PUT("/:articleID", articleService.Update)
	adminArticles.DELETE("/:articleID", articleService.Delete)

	return router
}
