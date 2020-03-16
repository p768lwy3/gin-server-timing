package servertiming

import (
	"github.com/gin-gonic/gin"
)

const contextKey = "Timing-Context"

func NewContext(c *gin.Context, h *Header) {
	c.Set(contextKey, h)
}

func FromContext(c *gin.Context) *Header {
	h, _ := c.MustGet(contextKey).(*Header)
	return h
}
