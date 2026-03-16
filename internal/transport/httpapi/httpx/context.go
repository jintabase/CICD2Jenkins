package httpx

import (
	"context"

	"cicd2jenkins/internal/domain"
)

type contextKey string

const actorContextKey contextKey = "actor"

func WithActor(ctx context.Context, actor domain.Actor) context.Context {
	return context.WithValue(ctx, actorContextKey, actor)
}

func ActorFromContext(ctx context.Context) (domain.Actor, bool) {
	actor, ok := ctx.Value(actorContextKey).(domain.Actor)
	return actor, ok
}
