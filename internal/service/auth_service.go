package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/logic"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type AuthService struct {
	auth *logic.AuthLogic
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(auth *logic.AuthLogic) *AuthService {
	return &AuthService{auth: auth}
}

func (s *AuthService) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := s.auth.Login(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		switch err {
		case apperrors.ErrBadRequest:
			httpx.WriteError(c, http.StatusBadRequest, err.Error())
		case apperrors.ErrInvalidCredentials:
			httpx.WriteError(c, http.StatusUnauthorized, err.Error())
		default:
			httpx.WriteError(c, http.StatusInternalServerError, "login failed")
		}
		return
	}

	httpx.WriteJSON(c, http.StatusOK, result)
}

func (s *AuthService) Me(c *gin.Context) {
	actor, ok := httpx.ActorFromContext(c)
	if !ok {
		httpx.WriteError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	httpx.WriteJSON(c, http.StatusOK, actor)
}
