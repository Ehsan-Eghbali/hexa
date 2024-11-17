package middleware

import (
	"github.com/gin-gonic/gin"
)

// SomeMiddleWare creates a new context with timeout for each request
func SomeMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
