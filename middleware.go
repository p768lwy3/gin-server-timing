package servertiming

import (
	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a variable of Header to time the server process
		h := Header{}

		// Set header to the context for the handler to call
		NewContext(c, &h)

		// Next handler
		c.Next()
	}
}

func WriteHeader(c *gin.Context) {
	// Get timing header from context
	h, _ := c.MustGet(contextKey).(*Header)

	// Write context to the header
	c.Header(HeaderKey, h.String())
}
