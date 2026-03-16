package middleware

import (
	"net/http"
	"strings"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
	"cicd2jenkins/internal/service"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type AuthMiddleware struct {
	auth *service.AuthService
}

func NewAuthMiddleware(auth *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if header == "" {
			httpx.WriteError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))
		actor, err := m.auth.ParseToken(token)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		next.ServeHTTP(w, r.WithContext(httpx.WithActor(r.Context(), actor)))
	})
}

func RequireRoles(roles ...domain.Role) func(http.Handler) http.Handler {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actor, ok := httpx.ActorFromContext(r.Context())
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
				return
			}

			if _, exists := allowed[actor.Role]; !exists {
				httpx.WriteError(w, http.StatusForbidden, apperrors.ErrForbidden.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
