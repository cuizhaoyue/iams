package middleware

import (
	"github.com/cuizhaoyue/iams/pkg/log"
	"github.com/gin-gonic/gin"
)

const UsernameKey = "username"

// Context 是一个中间件，它插入一些通用字段到gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(log.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
