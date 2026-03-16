package handler

import (
	"encoding/json"
	"net/http"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type AuthHandler struct {
	auth *service.AuthService
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.auth.Login(r.Context(), request.Username, request.Password)
	if err != nil {
		switch err {
		case apperrors.ErrBadRequest:
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
		case apperrors.ErrInvalidCredentials:
			httpx.WriteError(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	actor, ok := httpx.ActorFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, actor)
}
