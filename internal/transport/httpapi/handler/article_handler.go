package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type ArticleHandler struct {
	articles *service.ArticleService
}

func NewArticleHandler(articles *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articles: articles}
}

func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	articles, err := h.articles.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "list articles failed")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	article, err := h.articles.GetByID(r.Context(), chi.URLParam(r, "articleID"))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			httpx.WriteError(w, http.StatusNotFound, apperrors.ErrNotFound.Error())
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "get article failed")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor, ok := httpx.ActorFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input service.UpsertArticleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.articles.Create(r.Context(), actor, input)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "create article failed")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, article)
}

func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	actor, ok := httpx.ActorFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input service.UpsertArticleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.articles.Update(r.Context(), actor, chi.URLParam(r, "articleID"), input)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, apperrors.ErrNotFound.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "update article failed")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actor, ok := httpx.ActorFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.articles.Delete(r.Context(), actor, chi.URLParam(r, "articleID")); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrForbidden):
			httpx.WriteError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
		case errors.Is(err, apperrors.ErrBadRequest):
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, apperrors.ErrNotFound):
			httpx.WriteError(w, http.StatusNotFound, apperrors.ErrNotFound.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "delete article failed")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "article deleted",
	})
}
