package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	// XRequestIDKey 定义 X-Request-ID key 字符串.
	XRequestIDKey = "X-Request-ID"
)

// RequestID 是一个中间件，它插入'X-Request-ID'到上下文和每个请求的request/response的header中
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查传入的header是否有请求id，如果有则使用它
		rid := c.GetHeader(XRequestIDKey)

		// 如果没有请求id则创建一个uuid插入到header中和Context中
		if rid == "" {
			rid = uuid.NewV4().String()
			c.Request.Header.Set(XRequestIDKey, rid)
			c.Set(XRequestIDKey, rid)
		}

		// response中也要写入请求id
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}

// GetLoggerConfig 返回gin.LoggerConfig, 并指定日志写入的io.Writer和gin.LogFormatter.
// By default gin.DefaultWriter = os.Stdout
// reference: https://github.com/gin-gonic/gin#custom-log-format
func GetLoggerConfig(
	formatter gin.LogFormatter,
	output io.Writer,
	skipPaths []string,
) gin.LoggerConfig {
	if formatter == nil {
		formatter = GetDefaultLogFormatterWithRequestID()
	}

	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

// GetDefaultLogFormatterWithRequestID 返回带有'RequestID' 的gin.LogFormatter.
func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency -= param.Latency % time.Second
		}

		return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
			// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.ClientIP,
			param.Latency,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(XRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(XRequestIDKey)
}
