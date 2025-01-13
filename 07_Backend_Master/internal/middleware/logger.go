package middleware

import (
	"bytes"
	"io"
	"time"

	"backend-master/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// bodyLogWriter is a custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger middleware for request logging
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Set request ID in context
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID if available
		userID, _ := c.Get("user_id")

		// Log request details
		fields := map[string]interface{}{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        path,
			"query":       raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"duration":   duration.String(),
			"status":     c.Writer.Status(),
		}

		if userID != nil {
			fields["user_id"] = userID
		}

		// Log request body for non-GET requests (limited size)
		if c.Request.Method != "GET" && len(requestBody) > 0 {
			if len(requestBody) > 1024 {
				fields["request_body"] = string(requestBody[:1024]) + "..."
			} else {
				fields["request_body"] = string(requestBody)
			}
		}

		// Log response body (limited size)
		responseBody := blw.body.String()
		if len(responseBody) > 1024 {
			fields["response_body"] = responseBody[:1024] + "..."
		} else {
			fields["response_body"] = responseBody
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			log.Error("Server error", nil, fields)
		} else if c.Writer.Status() >= 400 {
			log.Warn("Client error", fields)
		} else {
			log.Info("Request processed", fields)
		}
	}
}

// RequestID middleware adds request ID to the response headers
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("request_id", requestID)
		}
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
