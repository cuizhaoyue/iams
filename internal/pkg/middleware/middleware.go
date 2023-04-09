package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  gin.Recovery(),
		"secure":    Secure,
		"options":   Options,
		"nocache":   NoCache,
		"cors":      Cors(),
		"requestid": RequestID(),
		"logger":    Logger(),
		"dump":      gindump.Dump(),
	}
}

// NoCache 是一个中间件函数，它附加标头以防止客户端缓存 HTTP 响应。
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options 是一个中间件函数，它为options请求附加标头并中止然后退出中间件链并结束请求.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Secure 是一个中间件，它添加了安全和资源访问相关的头.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")

	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}
