// Package middleware provides HTTP middleware functions for the NWC API
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoggingMiddleware logs information about each request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Time request start
		startTime := time.Now()

		// Process request
		c.Next()

		// After request processing
		duration := time.Since(startTime)
		
		// Log request details
		c.Request.Header.Del("X-API-Key") // Don't log the API key
		
		// Log the request
		gin.DefaultWriter.Write([]byte("[" + startTime.Format(time.RFC3339) + "] " +
			c.ClientIP() + " " +
			c.Request.Method + " " +
			c.Request.URL.Path + " " +
			c.Request.Proto + " " +
			http.StatusText(c.Writer.Status()) + " " +
			duration.String() + "\n"))
	}
}
