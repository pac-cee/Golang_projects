package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"backend-master/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Recovery middleware handles panics and logs them
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]

				// Get request details
				httpRequest := c.Request
				headers := make(map[string][]string)
				for k, v := range httpRequest.Header {
					headers[k] = v
				}

				// Read request body
				var requestBody []byte
				if httpRequest.Body != nil {
					requestBody, _ = io.ReadAll(httpRequest.Body)
					httpRequest.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				}

				// Log error with details
				fields := map[string]interface{}{
					"error":        fmt.Sprint(err),
					"request_id":   c.GetString("request_id"),
					"method":       httpRequest.Method,
					"path":         httpRequest.URL.Path,
					"query":        httpRequest.URL.RawQuery,
					"ip":          c.ClientIP(),
					"user_agent":  httpRequest.UserAgent(),
					"stack_trace": string(stack),
					"headers":     headers,
				}

				if len(requestBody) > 0 {
					if len(requestBody) > 1024 {
						fields["request_body"] = string(requestBody[:1024]) + "..."
					} else {
						fields["request_body"] = string(requestBody)
					}
				}

				log.Error("Recovery from panic", fmt.Errorf("%v", err), fields)

				// Return error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}

// CustomRecovery middleware allows custom error handling
func CustomRecovery(handler gin.RecoveryFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handler(c, err)
			}
		}()
		c.Next()
	}
}
