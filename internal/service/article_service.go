package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/logic"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type ArticleService struct {
	articles *logic.ArticleLogic
}

func NewArticleService(articles *logic.ArticleLogic) *ArticleService {
	return &ArticleService{articles: articles}
}

func (s *ArticleService) List(c *gin.Context) {
	articles, err := s.articles.List(c.Request.Context())
	if err != nil {
		httpx.WriteError(c, http.StatusInternalServerError, "list articles failed")
		return
	}
	httpx.WriteJSON(c, http.StatusOK, articles)
}

func (s *ArticleService) GetByID(c *gin.Context) {
	article, err := s.articles.GetByID(c.Request.Context(), c.Param("articleID"))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			httpx.WriteError(c, http.StatusNotFound, apperrors.ErrNotFound.Error())
			return
		}
		httpx.WriteError(c, http.StatusInternalServerError, "get article failed")
		return
	}
	httpx.WriteJSON(c, http.StatusOK, article)
}

func (s *ArticleService) Create(c *gin.Context) {
	actor, ok := httpx.ActorFromContext(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input logic.UpsertArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := s.articles.Create(c.Request.Context(), actor, input)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(c, http.StatusBadRequest, err.Error())
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "create article failed")
		}
		return
	}

	httpx.WriteJSON(c, http.StatusCreated, article)
}

func (s *ArticleService) Update(c *gin.Context) {
	actor, ok := httpx.ActorFromContext(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input logic.UpsertArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := s.articles.Update(c.Request.Context(), actor, c.Param("articleID"), input)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			httpx.WriteError(c, http.StatusNotFound, apperrors.ErrNotFound.Error())
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "update article failed")
		}
		return
	}

	httpx.WriteJSON(c, http.StatusOK, article)
}

func (s *ArticleService) Delete(c *gin.Context) {
	actor, ok := httpx.ActorFromContext(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := s.articles.Delete(c.Request.Context(), actor, c.Param("articleID")); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(c, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			httpx.WriteError(c, http.StatusNotFound, apperrors.ErrNotFound.Error())
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "delete article failed")
		}
		return
	}

	httpx.WriteJSON(c, http.StatusOK, map[string]string{
		"message": "article deleted",
	})
}
