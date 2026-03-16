package httpx

import (
	"github.com/gin-gonic/gin"
)

func WriteJSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

func WriteError(c *gin.Context, status int, message string) {
	WriteJSON(c, status, map[string]string{
		"error": message,
	})
}
