package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/logic"
	"cicd2jenkins/internal/model"
	"cicd2jenkins/internal/transport/httpapi/httpx"
)

type AuthMiddleware struct {
	auth *logic.AuthLogic
}

func NewAuthMiddleware(auth *logic.AuthLogic) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			httpx.WriteError(c, 401, "missing Authorization header")
			c.Abort()
			return
		}

		scheme, token, found := strings.Cut(header, " ")
		if !found || !strings.EqualFold(strings.TrimSpace(scheme), "Bearer") {
			httpx.WriteError(c, 401, "invalid Authorization header")
			c.Abort()
			return
		}

		actor, err := m.auth.ParseToken(strings.TrimSpace(token))
		if err != nil {
			httpx.WriteError(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		httpx.WithActor(c, actor)
		c.Next()
	}
}

func RequireRoles(roles ...model.Role) gin.HandlerFunc {
	allowed := make(map[model.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		actor, ok := httpx.ActorFromContext(c)
		if !ok {
			httpx.WriteError(c, 401, apperrors.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if _, exists := allowed[actor.Role]; !exists {
			httpx.WriteError(c, 403, apperrors.ErrForbidden.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
