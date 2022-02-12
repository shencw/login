package middleware

import (
	"github.com/gin-gonic/gin"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requestID", c.GetString(XRequestIDKey))
		c.Next()
	}
}
