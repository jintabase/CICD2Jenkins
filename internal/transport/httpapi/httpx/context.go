package httpx

import (
	"github.com/gin-gonic/gin"

	"cicd2jenkins/internal/model"
)

const actorContextKey = "actor"

func WithActor(c *gin.Context, actor model.Actor) {
	c.Set(actorContextKey, actor)
}

func ActorFromContext(c *gin.Context) (model.Actor, bool) {
	value, ok := c.Get(actorContextKey)
	if !ok {
		return model.Actor{}, false
	}

	actor, ok := value.(model.Actor)
	return actor, ok
}
