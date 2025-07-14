package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

// GinLogger gin logger is a middleware that logs the request info
// GinLogger 是一个中间件，用于记录请求信息
func GinLogger(logger *Logger) gin.HandlerFunc {
	//record request info
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"error":      param.ErrorMessage,
			"body_size":  param.BodySize,
			"user_agent": param.Request.UserAgent(),
			"request_id": param.Request.Header.Get("X-Request-Id"),
		}).Info("HTTP Request")

		return ""
	})

}

// GinRecovery gin recovery is a middleware that recovers from any panics and logs the error
// GinRecovery 是一个中间件，用于恢复任何恐慌并记录错误
func GinRecovery(logger *Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.WithFields(logrus.Fields{
			"error":      recovered,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"request_id": c.Request.Header.Get("X-Request-Id"),
		}).Error("HTTP Request")

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RequestID is a middleware that generates a request ID and stores it in the context
// RequestID 是一个中间件，用于生成请求ID并存储在上下文中
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = GenerateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// GenerateRequestID 生成请求ID
func GenerateRequestID() string {
	//time.Now() method is used to get the current time that will return a time.Time type value like 2006-01-02 15:04:05
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// rand.Intn() method is used to generate a random number
		//charset is used to generate a random string
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
